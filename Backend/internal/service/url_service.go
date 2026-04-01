package service

import (
	"Server/internal/entities/stats"
	"Server/internal/entities/url"
	"Server/internal/storage"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type analyticsEvent struct {
	shortURLID uint64
	ipAddress  string
	userAgent  string
}

type urlService struct {
	longURLRepo  storage.LongURLRepo
	shortURLRepo storage.ShortURLRepo
	statsRepo    storage.StatsRepo
	cache        storage.URLCache
	analyticsChan chan analyticsEvent
}

func NewURLService(longURLRepo storage.LongURLRepo, shortURLRepo storage.ShortURLRepo, statsRepo storage.StatsRepo, cache storage.URLCache) URLService {
	s := &urlService{
		longURLRepo:   longURLRepo,
		shortURLRepo:  shortURLRepo,
		statsRepo:     statsRepo,
		cache:         cache,
		analyticsChan: make(chan analyticsEvent, 1000),
	}

	// Start worker pool for analytics
	for i := 0; i < 5; i++ {
		go s.analyticsWorker()
	}

	return s
}

func (s *urlService) ShortenURL(ctx context.Context, userID string, originalURL string, customCode string) (*url.ShortURL, *url.LongURL, error) {
	// 1. Find or create the long URL
	lu, err := s.longURLRepo.FindOrCreate(ctx, originalURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store long URL: %w", err)
	}

	// 2. Generate short code
	code := customCode
	if code == "" {
		code = s.generateShortCode(originalURL)
	}

	// 3. Create short URL record
	su := &url.ShortURL{
		UserID:    userID,
		LongURLID: lu.ID,
		Code:      code,
	}

	err = s.shortURLRepo.Create(ctx, su)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	// 4. Cache for fast redirect
	_ = s.cache.SetLongURL(ctx, code, lu.Url, nil)

	return su, lu, nil
}


func (s *urlService) GetLongURL(ctx context.Context, code string) (string, error) {
	// Cache-aside: check cache first
	cached, err := s.cache.GetLongURL(ctx, code)
	if err == nil && cached != "" {
		return cached, nil
	}

	// Fallback to DB
	su, err := s.shortURLRepo.GetByCode(ctx, code)
	if err != nil || su == nil {
		return "", fmt.Errorf("URL not found")
	}

	lu, err := s.longURLRepo.GetByID(ctx, su.LongURLID)
	if err != nil || lu == nil {
		return "", fmt.Errorf("long URL not found")
	}

	// Populate cache
	_ = s.cache.SetLongURL(ctx, code, lu.Url, nil)

	return lu.Url, nil
}

func (s *urlService) GetUserURLs(ctx context.Context, userID string, limit, offset int) ([]URLDetail, int64, error) {
	shortURLs, total, err := s.shortURLRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	details := make([]URLDetail, 0, len(shortURLs))
	for _, su := range shortURLs {
		lu, err := s.longURLRepo.GetByID(ctx, su.LongURLID)
		longURL := ""
		if err == nil && lu != nil {
			longURL = lu.Url
		}
		details = append(details, URLDetail{
			ShortURLID: su.ID,
			Code:       su.Code,
			LongURL:    longURL,
			CreatedAt:  su.CreatedAt.Format(time.RFC3339),
			ExpiresAt:  su.ExpiresAt.Format(time.RFC3339),
		})
	}

	return details, total, nil
}

func (s *urlService) GetURLDetail(ctx context.Context, userID string, shortURLID uint64) (*URLDetail, error) {
	su, err := s.shortURLRepo.GetByID(ctx, shortURLID)
	if err != nil || su == nil {
		return nil, fmt.Errorf("URL not found")
	}

	lu, err := s.longURLRepo.GetByID(ctx, su.LongURLID)
	if err != nil || lu == nil {
		return nil, fmt.Errorf("long URL not found")
	}

	return &URLDetail{
		ShortURLID: su.ID,
		Code:       su.Code,
		LongURL:    lu.Url,
		CreatedAt:  su.CreatedAt.Format(time.RFC3339),
		ExpiresAt:  su.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *urlService) DeleteURL(ctx context.Context, userID string, shortURLID uint64) error {
	su, err := s.shortURLRepo.GetByID(ctx, shortURLID)
	if err != nil || su == nil {
		return fmt.Errorf("URL not found")
	}

	err = s.shortURLRepo.Delete(ctx, userID, shortURLID)
	if err != nil {
		return err
	}

	_ = s.cache.DeleteURL(ctx, su.Code)
	return nil
}

func (s *urlService) GetAnalytics(ctx context.Context, userID string, shortURLID uint64) (*stats.UrlStats, error) {
	// Verify ownership
	su, err := s.shortURLRepo.GetByID(ctx, shortURLID)
	if err != nil || su == nil {
		return nil, fmt.Errorf("URL not found")
	}

	return s.statsRepo.GetStats(ctx, shortURLID)
}

func (s *urlService) RecordClick(ctx context.Context, code string, ipAddress string, userAgent string) {
	su, err := s.shortURLRepo.GetByCode(ctx, code)
	if err != nil || su == nil {
		return
	}

	select {
	case s.analyticsChan <- analyticsEvent{shortURLID: su.ID, ipAddress: ipAddress, userAgent: userAgent}:
	default:
		log.Printf("Analytics channel full, dropping click for %s", code)
	}
}

// --- Internal helpers ---

func (s *urlService) generateShortCode(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hasher.Sum(nil)
	bits := binary.BigEndian.Uint64(hash[:8])
	code := s.base62Encode(bits)
	if len(code) > 10 {
		code = code[:10]
	}
	return code
}

func (s *urlService) base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Chars[0])
	}
	res := ""
	n := new(big.Int).SetUint64(num)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	for n.Cmp(zero) > 0 {
		rem := new(big.Int)
		n.QuoRem(n, base, rem)
		res = string(base62Chars[rem.Int64()]) + res
	}
	return res
}

func (s *urlService) analyticsWorker() {
	for event := range s.analyticsChan {
		ctx := context.Background()
		// Record individual click
		click := &stats.Click{
			ShortURLID: event.shortURLID,
			IpAddress:  event.ipAddress,
			UserAgent:  event.userAgent,
		}
		if err := s.statsRepo.RecordClick(ctx, click); err != nil {
			log.Printf("Failed to record click: %v", err)
		}
		// Increment aggregate stats
		if err := s.statsRepo.IncrementStats(ctx, event.shortURLID); err != nil {
			log.Printf("Failed to increment stats: %v", err)
		}
	}
}

package service

import (
	"Server/internal/entities/auth"
	"Server/internal/entities/user"
	"Server/internal/storage"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      storage.UserRepo
	jwtSecret string
}

func NewUserService(repo storage.UserRepo, jwtSecret string) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, u *user.User) (*user.User, error) {
	// Check if user already exists
	existingByEmail, _ := s.repo.GetUserByEmail(ctx, u.Email)
	if existingByEmail != nil {
		return nil, errors.New("email already registered")
	}

	existingByUsername, _ := s.repo.GetUserByUsername(ctx, u.UserName)
	if existingByUsername != nil {
		return nil, errors.New("username already taken")
	}

	// Set UUID and timestamps
	u.ID = uuid.New().String()
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err = s.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	// Don't return password
	u.Password = ""
	return u, nil
}

func (s *userService) Login(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	// For demo purposes, we will treat any client_id/client_secret as valid 
	// unless specified otherwise.
	
	u, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil || u == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateToken(u.ID)
}

func (s *userService) RefreshToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	// In a real OAuth2 system, we would check the refresh token in the DB.
	// For this, we'll parse the refresh token and generate a new pair.
	
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	uid := claims["uid"].(string)
	return s.generateToken(uid)
}

func (s *userService) RevokeToken(ctx context.Context, token string) error {
	// Implement revocation logic (e.g. adding to a Redis blacklist)
	return nil 
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) generateToken(uid string) (*auth.TokenResponse, error) {
	// Access Token
	now := time.Now()
	exp := now.Add(time.Hour * 1)
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"type": "access",
		"iat":  now.Unix(),
		"exp":  exp.Unix(),
	})

	accessStr, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshExp := now.Add(time.Hour * 24 * 7)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  uid,
		"type": "refresh",
		"iat":  now.Unix(),
		"exp":  refreshExp.Unix(),
	})

	refreshStr, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &auth.TokenResponse{
		AccessToken:  accessStr,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshStr,
	}, nil
}

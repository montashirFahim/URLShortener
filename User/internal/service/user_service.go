package service

import (
	"User/internal/domain/user"
	"User/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      repository.UserRepo
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepo, jwtSecret []byte) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, u *user.User) (int64, error) {
	// Check if user already exists
	existing, err := s.repo.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("user already exists with this email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.Password = string(hashedPassword)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return s.repo.CreateUser(ctx, u)
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.repo.GetUserByEmail(ctx, username) // Assuming email is used as username for login or we can check both
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": u.UID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

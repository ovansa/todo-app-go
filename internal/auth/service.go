package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo-app/internal/errors"
	"todo-app/internal/model"
	"todo-app/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Service interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, authUser *model.AuthUser) (string, error)
	ParseToken(tokenString string) (*Claims, error)
	AuthMiddleware() gin.HandlerFunc
	GetPepper() string
}

type authService struct {
	jwtSecret     string
	jwtExpiration time.Duration
	pepper        string
	userRepo      repository.UserRepository
}

func NewAuthService(jwtSecret string, jwtExpiration time.Duration, pepper string, userRepo repository.UserRepository) *authService {
	return &authService{
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
		pepper:        pepper,
		userRepo:      userRepo,
	}
}

func (s *authService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	if err := user.HashPassword(s.pepper); err != nil {
		return nil, errors.NewInternalServerError()
	}
	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, authUser *model.AuthUser) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, authUser.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.ErrInvalidCredentials
	}

	if err := user.ComparePassword(authUser.Password, s.pepper); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	claims := &Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ParseToken(tokenString string) (*Claims, error) {
	// First verify the token is not empty
	if tokenString == "" {
		return nil, errors.NewAPIError(http.StatusBadRequest, "INVALID", "Empty token string")
	}

	// Parse with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return nil, err
	}

	// Verify claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewAPIError(http.StatusBadRequest, "INVALID", "invalid token claims")
}

func (s *authService) GetPepper() string {
	return s.pepper
}

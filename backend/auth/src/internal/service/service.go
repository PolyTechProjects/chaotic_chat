package service

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/PolyTechProjects/chaotic_chat/auth/src/internal/repository"
	"github.com/PolyTechProjects/chaotic_chat/auth/src/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
	jwtSecretKey   []byte
	keyFunc        func(token *jwt.Token) (interface{}, error)
}

func New(authRepository *repository.AuthRepository) *AuthService {
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	}
	return &AuthService{AuthRepository: authRepository, jwtSecretKey: jwtSecretKey, keyFunc: keyFunc}
}

func (s *AuthService) Register(login string, password string) (string, uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", uuid.Nil, err
	}
	user, err := models.New(login, string(hash))
	if err != nil {
		return "", uuid.Nil, err
	}
	err = s.AuthRepository.Save(user)
	if err != nil {
		return "", uuid.Nil, err
	}

	accessToken, err := s.generateAccessToken(user.Id)
	if err != nil {
		return "", uuid.Nil, err
	}
	slog.Info(fmt.Sprintf("User %v registered", user.Id))
	return accessToken, user.Id, nil
}

func (s *AuthService) Login(login string, password string) (string, error) {
	user, err := s.AuthRepository.FindByLogin(login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(password))
	if err != nil {
		return "", err
	}

	accessToken, err := s.generateAccessToken(user.Id)
	slog.Debug(fmt.Sprintf("accessToken: %v", accessToken))
	if err != nil {
		return "", err
	}
	slog.Info(fmt.Sprintf("User %v authenticated", user.Id))
	return accessToken, nil
}

func (s *AuthService) Authorize(accessToken string) (string, uuid.UUID, error) {
	var claims jwt.MapClaims
	_, err := jwt.ParseWithClaims(accessToken, &claims, s.keyFunc)
	if err != nil {
		return "", uuid.Nil, err
	}

	userId, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return "", uuid.Nil, err
	}
	slog.Debug(fmt.Sprintf("userId: %v", userId))

	user, err := s.AuthRepository.FindById(userId)
	if err != nil {
		return "", uuid.Nil, err
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		accessToken, err = s.generateAccessToken(user.Id)
		if err != nil {
			return "", uuid.Nil, err
		}
	}

	return accessToken, userId, nil
}

func (s *AuthService) ExtractUserId(tokenString string) (string, error) {
	var claims jwt.MapClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, s.keyFunc)
	if err != nil {
		return "", err
	}
	return claims["sub"].(string), nil
}

func (s *AuthService) generateAccessToken(userId uuid.UUID) (string, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	slog.Info(fmt.Sprintf("exp: %v", exp))
	payload := jwt.MapClaims{
		"sub": userId,
		"exp": exp,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(s.jwtSecretKey)
}

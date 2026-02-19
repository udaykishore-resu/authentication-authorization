package service

import (
	"auth-service/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler interface {
	Register(empType, username, password string) error
	Login(username, password string) (empID int, accessToken string, refreshToken string, err error)
}

type EmpRepository interface {
	CreateEmp(empType, username, passwordHash string) error
	GetEmpByEmpName(username string) (empID int, passwordHash string, err error)
}

type TokenRepository interface {
	StoreRefreshToken(empID int, token string, expires time.Duration) error
	GetEmpIDByRefreshToken(token string) (int, error)
	DeleteRefreshToken(token string) error
}

type AuthServiceImpl struct {
	empRepo   EmpRepository
	tokenRepo TokenRepository
}

func NewAuthService(empRepo EmpRepository, tokenRepo TokenRepository) AuthHandler {
	return &AuthServiceImpl{
		empRepo:   empRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthServiceImpl) Register(empType, username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.empRepo.CreateEmp(empType, username, string(hashed)); err != nil {
		return errors.New("failed to create employee")
	}

	return nil
}

func (s *AuthServiceImpl) Login(username, password string) (int, string, string, error) {
	empID, storedHash, err := s.empRepo.GetEmpByEmpName(username)
	if err != nil {
		return 0, "", "", errors.New("employee not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return 0, "", "", errors.New("invalid credentials")
	}

	refreshToken := uuid.New().String()
	accessToken, err := utils.GenerateJWT(empID)
	if err != nil {
		return 0, "", "", errors.New("failed to generate access token")
	}

	if err := s.tokenRepo.StoreRefreshToken(empID, refreshToken, 24*time.Hour*7); err != nil {
		return 0, "", "", errors.New("failed to store refresh token")
	}

	return empID, accessToken, refreshToken, nil
}

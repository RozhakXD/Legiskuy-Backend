package auth

import (
	"errors"
	"legiskuy-backend/internal/voter"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input *RegisterInput) (*UserResponse, error)
	Login(input *LoginInput) (string, error)
}

type service struct {
	repository Repository
	voterRepo  voter.Repository
}

func NewService(repo Repository, voterRepo voter.Repository) Service {
	return &service{
		repository: repo,
		voterRepo:  voterRepo,
	}
}

func (s *service) Register(input *RegisterInput) (*UserResponse, error) {
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.Username == "" {
		return nil, errors.New("username is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	existingUser, _ := s.repository.FindByUsername(input.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     input.Name,
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "pemilih",
	}

	newUser, err := s.repository.Create(user)
	if err != nil {
		return nil, err
	}

	voterData := &voter.Voter{
		Name:     newUser.Name,
		HasVoted: false,
	}
	s.voterRepo.Create(voterData)

	response := &UserResponse{
		ID:       newUser.ID,
		Name:     newUser.Name,
		Username: newUser.Username,
		Role:     newUser.Role,
	}
	return response, nil
}

func (s *service) Login(input *LoginInput) (string, error) {
	user, err := s.repository.FindByUsername(input.Username)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

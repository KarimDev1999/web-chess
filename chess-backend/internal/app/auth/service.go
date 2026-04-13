package auth

import (
	"context"

	"chess-backend/internal/domain/user"
	"chess-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo user.Repository
	tokenSvc *jwt.TokenService
}

func NewAuthService(userRepo user.Repository, tokenSvc *jwt.TokenService) *AuthService {
	return &AuthService{userRepo: userRepo, tokenSvc: tokenSvc}
}

func (s *AuthService) Register(ctx context.Context, cmd RegisterCommand) (*user.User, string, error) {
	existing, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrEmailAlreadyRegistered
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	u := user.NewUser(cmd.Email, string(hashed), cmd.Username)
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, "", err
	}
	token, err := s.tokenSvc.Generate(u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func (s *AuthService) Login(ctx context.Context, cmd LoginCommand) (*user.User, string, error) {
	u, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(cmd.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := s.tokenSvc.Generate(u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

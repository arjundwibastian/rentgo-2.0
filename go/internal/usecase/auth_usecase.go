package usecase

import (
	"context"
	"errors"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"uuid"

	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	repo domain.UserRepository
}

func NewAuthUseCase(repo domain.UserRepository) domain.AuthUseCase {
	return &authUseCase{repo: repo}
}

func (ac *authUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*domain.User, error) {
	if _, err := ac.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, domain.ErrEmailRegistered
	} else if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil , err
	}
	uuid := uuid.NewV7()
	
	user := &domain.User{
		ID		: uuid,
		FullName: req.FullName,
		PasswordHash: string(hashed),
		PhoneNumber : req.PhoneNumber,
		Address : req.Address,
		Role : domain.RoleUser,
		IsActive: true,
	}
	err = ac.repo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ac *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := ac.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}	
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, domain.ErrForbidden
	}

}
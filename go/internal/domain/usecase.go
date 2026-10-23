package domain

import (
	"context"
	"milestone-02/internal/dto"
)

type AuthUseCase interface {
	Register (ctx context.Context, req dto.RegisterRequest) (*User, error)
	Login (ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type UserUseCase interface {
	
	
}

type AdminUseCase interface {
	
}
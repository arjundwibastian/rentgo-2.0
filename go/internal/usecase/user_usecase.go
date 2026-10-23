package usecase

import (
	"milestone-02/internal/domain"
)

type userUseCase struct {
	repo    domain.UserRepository
	emailService domain.EmailService
}

func NewUserUseCase(repo domain.UserRepository, emailService domain.EmailService) domain.UserUseCase {
	return &userUseCase{repo: repo, emailService: emailService}
}

func (uc *userUseCase) GetProfile(userID int) (*domain.User, error) {
	user, err := uc.repo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}


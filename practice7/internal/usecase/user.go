package usecase

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	createdUser, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}
	sessionID := uuid.New().String()
	return createdUser, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.LoginUser(user)
	if err != nil {
		return "", fmt.Errorf("user from repo: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, user.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return token, nil
}

// GetMe — возвращает пользователя по userID из JWT (Задание 1)
func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := u.repo.GetUserByID(parsedID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

// PromoteUser — меняет роль пользователя на admin (Задание 2)
func (u *UserUseCase) PromoteUser(userID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.PromoteUser(userID)
	if err != nil {
		return nil, fmt.Errorf("promote user: %w", err)
	}
	return user, nil
}

package usecase

import (
	"context"
	"fmt"
	"strings"

	"golang/internal/repository"
	"golang/pkg/modules"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]modules.User, error)
	GetUserByID(ctx context.Context, id int) (*modules.User, error)
	CreateUser(ctx context.Context, u modules.User) (*modules.User, error)
	UpdateUser(ctx context.Context, id int, patch UpdateUserPatch) (*modules.User, error)
	DeleteUser(ctx context.Context, id int) (int64, error)
}

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(ctx context.Context) ([]modules.User, error) {
	return u.repo.GetUsers(ctx)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *UserUsecase) CreateUser(ctx context.Context, user modules.User) (*modules.User, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if user.Name == "" {

		return nil, fmt.Errorf("%w: name is required", modules.ErrInvalidInput)
	}
	if user.Email == "" || !strings.Contains(user.Email, "@") {
		return nil, fmt.Errorf("%w: valid email is required", modules.ErrInvalidInput)
	}
	if user.Age < 0 {
		return nil, fmt.Errorf("%w: age must be >= 0", modules.ErrInvalidInput)
	}

	return u.repo.CreateUser(ctx, user)
}

type UpdateUserPatch struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id int, patch UpdateUserPatch) (*modules.User, error) {
	current, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated := *current

	if patch.Name != nil {
		name := strings.TrimSpace(*patch.Name)
		if name == "" {
			return nil, fmt.Errorf("%w: name cannot be empty", modules.ErrInvalidInput)
		}
		updated.Name = name
	}
	if patch.Email != nil {
		email := strings.TrimSpace(*patch.Email)
		if email == "" || !strings.Contains(email, "@") {
			return nil, fmt.Errorf("%w: valid email is required", modules.ErrInvalidInput)
		}
		updated.Email = email
	}
	if patch.Age != nil {
		if *patch.Age < 0 {
			return nil, fmt.Errorf("%w: age must be >= 0", modules.ErrInvalidInput)
		}
		updated.Age = *patch.Age
	}

	if err := u.repo.UpdateUser(ctx, id, updated); err != nil {
		return nil, err
	}

	return u.repo.GetUserByID(ctx, id)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int) (int64, error) {
	return u.repo.DeleteUser(ctx, id)
}

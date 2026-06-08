package err

import (
	"context"
	"errors"
	"fmt"
)

type userRepo interface {
	GetUser(ctx context.Context, email string) (*User, error)
	AddUser(ctx context.Context, user *User) error
}
type UserService struct {
	repo userRepo
}

func NewUserService(repo userRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, email, userName string) error {
	user := &User{
		email:    email,
		userName: userName,
	}

	usr, err := u.repo.GetUser(ctx, email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("userService:CreateUser:%w", err)
	}
	if usr != nil {
		return fmt.Errorf("userService:CreateUsesr:%w", ErrEmailAlredyExist)
	}
	if err := u.repo.AddUser(ctx, user); err != nil {
		return fmt.Errorf("userService:CreateUser:%w", err)
	}
	return nil
}

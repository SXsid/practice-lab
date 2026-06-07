package err

import "fmt"

type userRepo interface {
	GetUser(email string) (*User, error)
	AddUser(user *User) error
}
type UserService struct {
	repo userRepo
}

func NewUserService(repo userRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(email, userName string) error {
	user := &User{
		email:    email,
		userName: userName,
	}
	if err := u.repo.AddUser(user); err != nil {
		return fmt.Errorf("userService:CreateUser:%w", err)
	}
	return nil
}

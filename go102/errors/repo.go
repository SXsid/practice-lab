package err

import "fmt"

type UserRepo struct {
	data map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*User{},
	}
}

func (u *UserRepo) AddUser(user *User) error {
	if _, ok := u.data[user.email]; ok {
		return fmt.Errorf("userRepo:AddUser:%w", ErrEmailAlredyExist)
	}
	u.data[user.email] = user
	return nil
}

func (u *UserRepo) GetUser(email string) (*User, error) {
	user, ok := u.data[email]
	if !ok {
		return nil, fmt.Errorf("userRepo:GetUser:%w", ErrNotFound)
	}
	return user, nil
}

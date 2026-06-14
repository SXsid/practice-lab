package err

import (
	"context"
	"fmt"
	"sync"
)

type UserRepo struct {
	mu   sync.RWMutex
	data map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*User{},
	}
}

func (u *UserRepo) AddUser(ctx context.Context, user *User) error {
	u.mu.RLock()
	defer u.mu.RUnlock()
	u.data[user.email] = user
	return nil
}

func (u *UserRepo) GetUser(ctx context.Context, email string) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	user, ok := u.data[email]
	// in real if sql.rwo ==0 we wrap with out system eroro
	// so system talk with common eror type and last we can decide with status code need ot send to the forntend
	if !ok {
		return nil, fmt.Errorf("userRepo:GetUser:%w", ErrNotFound)
	}
	return user, nil
}

package err

import (
	"context"
	"fmt"
	"time"
)

type UserRepo struct {
	data map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*User{},
	}
}

// TODO: how to write custoem cacellbe context funtion whcih repect it
// best pracies
// may be some youtube video later
func (u *UserRepo) AddUser(ctx context.Context, user *User) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("userRepo:Addusesr:%w", ErrTimeout)
		default:

			time.Sleep(4 * time.Second)
			u.data[user.email] = user
			return nil
		}
	}
}

func (u *UserRepo) GetUser(ctx context.Context, email string) (*User, error) {
	user, ok := u.data[email]
	// in real if sql.rwo ==0 we wrap with out system eroro
	// so system talk with common eror type and last we can decide with status code need ot send to the forntend
	if !ok {
		return nil, fmt.Errorf("userRepo:GetUser:%w", ErrNotFound)
	}
	return user, nil
}

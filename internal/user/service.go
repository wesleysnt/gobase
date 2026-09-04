package user

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/wesleysnt/gobase/internal/httputil"
)

var validate = validator.New()

// UserRepository defines the data access interface for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", httputil.ErrValidation, err)
	}

	user := &User{Name: req.Name, Email: req.Email}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) List(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		r := u.ToResponse()
		resp[i] = &r
	}
	return resp, nil
}

func (s *UserService) Update(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", httputil.ErrValidation, err)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

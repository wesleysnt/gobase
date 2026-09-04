package user

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/wesleysnt/gobase/internal/httputil"
)

type stubUserRepo struct {
	users  map[string]*User
	nextID int
}

func newStubRepo() *stubUserRepo {
	return &stubUserRepo{users: make(map[string]*User), nextID: 1}
}

func (r *stubUserRepo) Create(ctx context.Context, u *User) error {
	u.ID = fmt.Sprintf("generated-id-%d", r.nextID)
	r.nextID++
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	// Simulate email uniqueness check
	for _, existing := range r.users {
		if existing.Email == u.Email {
			return httputil.ErrConflict
		}
	}
	r.users[u.ID] = u
	return nil
}

func (r *stubUserRepo) FindByID(ctx context.Context, id string) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, httputil.ErrNotFound
	}
	return u, nil
}

func (r *stubUserRepo) List(ctx context.Context) ([]*User, error) {
	result := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *stubUserRepo) Update(ctx context.Context, u *User) error {
	if _, ok := r.users[u.ID]; !ok {
		return httputil.ErrNotFound
	}
	u.UpdatedAt = time.Now()
	r.users[u.ID] = u
	return nil
}

func (r *stubUserRepo) Delete(ctx context.Context, id string) error {
	if _, ok := r.users[id]; !ok {
		return httputil.ErrNotFound
	}
	delete(r.users, id)
	return nil
}

func TestServiceCreate(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	req := CreateUserRequest{Name: "Alice", Email: "alice@example.com"}
	user, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.ID == "" {
		t.Error("Create() user.ID is empty")
	}
	if user.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", user.Name)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr bool
	}{
		{"missing name", CreateUserRequest{Name: "", Email: "a@b.com"}, true},
		{"missing email", CreateUserRequest{Name: "Bob", Email: ""}, true},
		{"short name", CreateUserRequest{Name: "AB", Email: "a@b.com"}, true},
		{"invalid email", CreateUserRequest{Name: "Charlie", Email: "not-an-email"}, true},
		{"valid", CreateUserRequest{Name: "Diana", Email: "diana@example.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tt.req)
			if tt.wantErr && err == nil {
				t.Error("Create() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Create() unexpected error: %v", err)
			}
		})
	}
}

func TestServiceGetByID(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	created, _ := svc.Create(context.Background(), CreateUserRequest{Name: "Eve", Email: "eve@example.com"})

	user, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if user.Name != "Eve" {
		t.Errorf("Name = %q, want Eve", user.Name)
	}
}

func TestServiceGetByIDNotFound(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	_, err := svc.GetByID(context.Background(), "nonexistent")
	if !errors.Is(err, httputil.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
}

func TestServiceList(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	svc.Create(context.Background(), CreateUserRequest{Name: "Alice", Email: "alice@example.com"})
	svc.Create(context.Background(), CreateUserRequest{Name: "Bob", Email: "bob@example.com"})

	users, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}
}

func TestServiceUpdate(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	created, _ := svc.Create(context.Background(), CreateUserRequest{Name: "Eve", Email: "eve@example.com"})

	newName := "Eve Updated"
	updated, err := svc.Update(context.Background(), created.ID, UpdateUserRequest{Name: &newName})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != "Eve Updated" {
		t.Errorf("Name = %q, want Eve Updated", updated.Name)
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	name := "Nobody"
	_, err := svc.Update(context.Background(), "nonexistent", UpdateUserRequest{Name: &name})
	if !errors.Is(err, httputil.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestServiceDelete(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	created, _ := svc.Create(context.Background(), CreateUserRequest{Name: "Charlie", Email: "charlie@example.com"})

	err := svc.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if !errors.Is(err, httputil.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	err := svc.Delete(context.Background(), "nonexistent")
	if !errors.Is(err, httputil.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

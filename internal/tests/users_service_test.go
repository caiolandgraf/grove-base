package tests

import (
	"errors"
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"github.com/caiolandgraf/grove-base/internal/modules/auth"
	"github.com/caiolandgraf/grove-base/internal/modules/users"
)

func TestUsersService(t *testing.T) {
	s := gest.Describe("UsersService")

	s.It("should create a user with a hashed password", func(t *gest.T) {
		svc := users.NewService(newMockUserStore())

		res, err := svc.CreateUser(&users.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		})

		t.Expect(err).ToBe(nil)
		t.Expect(res.Name).ToBe("Test User")
		t.Expect(res.Email).ToBe("test@example.com")
		t.Expect(res.ID).Not().ToBe("")
	})

	s.It("should prevent duplicate emails", func(t *gest.T) {
		svc := users.NewService(newMockUserStore())
		req := &users.CreateUserRequest{
			Name:     "Test User",
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		_, err := svc.CreateUser(req)
		t.Expect(err).ToBe(nil)

		_, err = svc.CreateUser(req)
		t.Expect(err).ToBe(users.ErrEmailAlreadyExists)
	})

	s.Run(t)
}

func TestAuthService(t *testing.T) {
	s := gest.Describe("AuthService")

	s.It("should authenticate users with correct credentials", func(t *gest.T) {
		store := newMockUserStore()
		userSvc := users.NewService(store)
		authSvc := auth.NewService(store)

		_, err := userSvc.CreateUser(&users.CreateUserRequest{
			Name:     "Auth User",
			Email:    "auth@example.com",
			Password: "securepassword",
		})
		t.Expect(err).ToBe(nil)

		user, err := authSvc.Login("auth@example.com", "securepassword")
		t.Expect(err).ToBe(nil)
		t.Expect(user.Email).ToBe("auth@example.com")

		_, err = authSvc.Login("auth@example.com", "wrongpassword")
		t.Expect(err).Not().ToBe(nil)
	})

	s.It("should reject login for unknown users", func(t *gest.T) {
		authSvc := auth.NewService(newMockUserStore())

		_, err := authSvc.Login("missing@example.com", "password")
		t.Expect(err).Not().ToBe(nil)
	})

	s.Run(t)
}

// mockUserStore is shared with users_service_test.go in this package.
type mockUserStore struct {
	users map[string]*users.User
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{
		users: make(map[string]*users.User),
	}
}

func (m *mockUserStore) Create(user *users.User) error {
	if user.ID == "" {
		user.ID = "mock-uuid-" + user.Email
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserStore) Update(user *users.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserStore) Delete(id any) error {
	idStr, ok := id.(string)
	if !ok {
		return errors.New("invalid id type")
	}
	if _, ok := m.users[idStr]; !ok {
		return errors.New("not found")
	}
	delete(m.users, idStr)
	return nil
}

func (m *mockUserStore) FindByID(id string) (*users.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, users.ErrNotFound
	}
	return user, nil
}

func (m *mockUserStore) FindByEmail(email string) (*users.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, users.ErrNotFound
}

func (m *mockUserStore) FindAll(page, perPage int) ([]users.User, int64, error) {
	list := make([]users.User, 0, len(m.users))
	for _, user := range m.users {
		list = append(list, *user)
	}
	return list, int64(len(list)), nil
}

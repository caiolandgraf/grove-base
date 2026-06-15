package tests

import (
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"github.com/caiolandgraf/grove-base/internal/modules/auth"
	"github.com/caiolandgraf/grove-base/internal/modules/users"
	"github.com/caiolandgraf/grove-base/internal/tests/testutil"
)

func TestUsersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupDB(t)

	s := gest.Describe("UsersIntegration")

	s.It("should persist and retrieve a user from PostgreSQL", func(t *gest.T) {
		t.Expect(testutil.TruncateUsers(db)).ToBe(nil)

		svc := users.NewService(users.Users(db))

		created, err := svc.CreateUser(&users.CreateUserRequest{
			Name:     "Integration User",
			Email:    "integration@example.com",
			Password: "password123",
		})
		t.Expect(err).ToBe(nil)
		t.Expect(created.ID).Not().ToBe("")

		found, err := svc.GetUserByID(created.ID)
		t.Expect(err).ToBe(nil)
		t.Expect(found.Email).ToBe("integration@example.com")
	})

	s.It("should enforce unique emails in the database", func(t *gest.T) {
		t.Expect(testutil.TruncateUsers(db)).ToBe(nil)

		svc := users.NewService(users.Users(db))
		req := &users.CreateUserRequest{
			Name:     "User One",
			Email:    "unique@example.com",
			Password: "password123",
		}

		_, err := svc.CreateUser(req)
		t.Expect(err).ToBe(nil)

		_, err = svc.CreateUser(req)
		t.Expect(err).ToBe(users.ErrEmailAlreadyExists)
	})

	s.Run(t)
}

func TestAuthIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutil.SetupDB(t)

	s := gest.Describe("AuthIntegration")

	s.It("should login with credentials stored in PostgreSQL", func(t *gest.T) {
		t.Expect(testutil.TruncateUsers(db)).ToBe(nil)

		repo := users.Users(db)
		userSvc := users.NewService(repo)
		authSvc := auth.NewService(repo)

		_, err := userSvc.CreateUser(&users.CreateUserRequest{
			Name:     "Auth Integration",
			Email:    "auth-integration@example.com",
			Password: "securepassword",
		})
		t.Expect(err).ToBe(nil)

		loggedIn, err := authSvc.Login("auth-integration@example.com", "securepassword")
		t.Expect(err).ToBe(nil)
		t.Expect(loggedIn.Email).ToBe("auth-integration@example.com")

		_, err = authSvc.Login("auth-integration@example.com", "wrongpassword")
		t.Expect(err).Not().ToBe(nil)
	})

	s.Run(t)
}

package tests

import (
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"github.com/caiolandgraf/grove-base/internal/models"
)

func TestUser(t *testing.T) {
	s := gest.Describe("User")

	s.It("should have valid fields", func(t *gest.T) {
		user := models.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		t.Expect(user.Name).ToBe("John Doe")
		t.Expect(user.Email).ToBe("john@example.com")
	})

	s.It("should have empty ID before persisting", func(t *gest.T) {
		user := models.User{}

		t.Expect(user.ID).ToBe("")
	})

	s.It("table name should be users", func(t *gest.T) {
		user := models.User{}

		t.Expect(user.TableName()).ToBe("users")
	})

	s.Run(t)
}

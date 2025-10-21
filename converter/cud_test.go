package converter

import (
	"context"
	"testing"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/testutils"
)

var repo gormen.Repository[testutils.User]

func TestMain(m *testing.M) {
	repo = NewRepository[testutils.UserDB, *testutils.UserDB](testutils.SetupTestDB())
	m.Run()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		user := &testutils.User{
			Username: "jdoe",
			Password: "1234",
			Person: testutils.Person{
				Name:  "John Doe",
				Email: "jdoe@mail.com",
			},
		}

		err := repo.Create(ctx, user)
		if err != nil || user.ID == 0 {
			t.Fatalf("Error creating User %v", err)
		}
	})

	t.Run("CreateAll", func(t *testing.T) {
		user := testutils.User{
			Username: "batch1",
			Password: "1234",
			Person: testutils.Person{
				Name:  "Batch 1",
				Email: "b1@mail.com",
			},
		}
		user2 := testutils.User{
			Username: "batch2",
			Password: "1234",
			Person: testutils.Person{
				Name:  "Batch 2",
				Email: "b2@mail.com",
			},
		}

		users := []testutils.User{user, user2}
		err := repo.CreateAll(ctx, &users, 2)

		if err != nil || users[0].ID == 0 || users[1].ID == 0 {
			t.Fatalf("Error creating users %v", err)
		}
	})
}

package std

import (
	"context"
	"testing"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/testutils"
)

var repo gormen.Repository[testutils.UserDB]

func TestMain(m *testing.M) {
	repo = NewRepository[testutils.UserDB](testutils.SetupTestDB())
	m.Run()
}

func TestCud(t *testing.T) {
	ctx := context.Background()

	user := testutils.UserDB{
		Username: "batch1",
		Password: "1234",
		Person: testutils.PersonDB{
			Name:  "Batch 1",
			Email: "b1@mail.com",
		},
	}
	user2 := testutils.UserDB{
		Username: "batch2",
		Password: "1234",
		Person: testutils.PersonDB{
			Name:  "Batch 2",
			Email: "b2@mail.com",
		},
	}
	user3 := testutils.UserDB{
		Username: "batch3",
		Password: "1234",
		Person: testutils.PersonDB{
			Name:  "Batch 3",
			Email: "b3@mail.com",
		},
	}
	users := []testutils.UserDB{user, user2, user3}

	t.Run("Std Create", func(t *testing.T) {
		user := &testutils.UserDB{
			Username: "jdoe",
			Password: "1234",
			Person: testutils.PersonDB{
				Name:  "John Doe",
				Email: "jdoe@mail.com",
			},
		}

		err := repo.Create(ctx, user)
		if err != nil || user.ID == 0 {
			t.Fatalf("Error creating User %v", err)
		}
	})

	t.Run("Std CreateAll", func(t *testing.T) {
		err := repo.CreateAll(ctx, &users, 2)
		if err != nil || users[0].ID == 0 || users[1].ID == 0 {
			t.Fatalf("Error creating users %v", err)
		}
	})

	t.Run("Std Save", func(t *testing.T) {
		user := users[0]
		user.Password = "1111"

		err := repo.Save(ctx, &user)
		if err != nil {
			t.Fatalf("Error saving User %v", err)
		}
	})

	t.Run("Std SaveAll", func(t *testing.T) {
		users[0].Password = "123"
		users[1].Password = "123"

		err := repo.SaveAll(ctx, users)
		if err != nil {
			t.Fatalf("Error saving User %v", err)
		}
	})

	t.Run("Std Delete", func(t *testing.T) {
		err := repo.Delete(ctx, &users[2])
		if err != nil {
			t.Fatalf("Error deleting User %v", err)
		}
	})
}

package main

import (
	"context"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/gormen/std"
	"github.com/javiorfo/gormen/where"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	var repo gormen.Repository[UserDB]
	repo = std.NewRepository[UserDB](SetupDB())

	user := &UserDB{
		Username: "jdoe",
		Password: "1234",
		Person: PersonDB{
			Name:  "John Doe",
			Email: "jdoe@mail.com",
		},
	}

	err := repo.Create(ctx, user)
	if err != nil || user.ID == 0 {
		log.Fatalf("Error creating User %v", err)
	}

	user1 := UserDB{
		Username: "batch1",
		Password: "1234",
		Person: PersonDB{
			Name:  "Batch 1",
			Email: "b1@mail.com",
		},
	}
	user2 := UserDB{
		Username: "batch2",
		Password: "1234",
		Person: PersonDB{
			Name:  "Batch 2",
			Email: "b2@mail.com",
		},
	}
	user3 := UserDB{
		Username: "batch3",
		Password: "1234",
		Person: PersonDB{
			Name:  "Batch 3",
			Email: "b3@mail.com",
		},
	}
	users := []UserDB{user1, user2, user3}

	err = repo.CreateAll(ctx, &users, 3)
	if err != nil {
		log.Fatalf("Error creating users %v", err)
	}

	opt, _ := repo.FindBy(ctx, gormen.NewWhere(where.Like("username", "%do%")).Build(), "Person")
	opt.Inspect(func(ud UserDB) {
		log.Printf("%+v", ud)
	})

	orders := []sort.Order{sort.NewOrder("username", sort.Descending)}
	users, err = repo.FindAllOrdered(ctx, orders, "Person")
	if err != nil {
		log.Fatalf("Error searching users %v", err)
	}
	log.Printf("%+v", users)
}

type UserDB struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Username string   `gorm:"not null"`
	Password string   `gorm:"not null"`
	PersonID uint     `gorm:"not null"`
	Person   PersonDB `gorm:"column:person_id;not null"`
}

func (udb UserDB) TableName() string {
	return "users"
}

type PersonDB struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (pdb PersonDB) TableName() string {
	return "people"
}

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&PersonDB{}, &UserDB{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

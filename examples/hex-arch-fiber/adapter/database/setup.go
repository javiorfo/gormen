package database

import (
	"hex-arch-fiber/adapter/database/entities"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&entities.PersonDB{}, &entities.UserDB{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	user := &entities.UserDB{
		Username: "batch1",
		Password: "1234",
		Enable:   true,
		Person: entities.PersonDB{
			Name:  "Batch 1",
			Email: "b1@mail.com",
		},
	}
	user2 := &entities.UserDB{
		Username: "batch2",
		Password: "1234",
		Enable:   true,
		Person: entities.PersonDB{
			Name:  "Batch 2",
			Email: "b2@mail.com",
		},
	}
	user3 := &entities.UserDB{
		Username: "batch3",
		Password: "1234",
		Enable:   true,
		Person: entities.PersonDB{
			Name:  "Batch 3",
			Email: "b3@mail.com",
		},
	}
	user4 := &entities.UserDB{
		Username: "batch4",
		Password: "1234",
		Enable:   false,
		Person: entities.PersonDB{
			Name:  "Batch 4",
			Email: "b4@mail.com",
		},
	}

	db.Create(user)
	db.Create(user2)
	db.Create(user3)
	db.Create(user4)

	return db
}

package testutils

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type UserDB struct {
	ID       uint     `gorm:"primaryKey;autoIncrement"`
	Username string   `gorm:"not null"`
	Password string   `gorm:"not null"`
	PersonID uint     `gorm:"not null"`
	Person   PersonDB `gorm:"column:person_id;not null"`
}

func (udb *UserDB) From(u User) {
	udb.ID = u.ID
	udb.Username = u.Username
	udb.Password = u.Password
	udb.Person.From(u.Person)
}

func (udb UserDB) Into() User {
	return User{
		ID:       udb.ID,
		Username: udb.Username,
		Password: udb.Password,
		Person:   udb.Person.Into(),
	}
}

type User struct {
	ID       uint
	Username string
	Password string
	Person   Person
}

type PersonDB struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

func (pdb *PersonDB) From(p Person) {
	pdb.ID = p.ID
	pdb.Name = p.Name
	pdb.Email = p.Email
}

func (pdb PersonDB) Into() Person {
	return Person{
		ID:    pdb.ID,
		Name:  pdb.Name,
		Email: pdb.Email,
	}
}

type Person struct {
	ID    uint
	Name  string
	Email string
}

func SetupTestDB() *gorm.DB {
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

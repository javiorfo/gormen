package main

import (
	"context"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/converter"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/gormen/where"
	"gorm.io/gorm"
)

type userRepo struct {
	gormen.Repository[User]
}

func NewUserRepo(db *gorm.DB) userRepo {
	return userRepo{
		Repository: converter.NewRepository[UserDB, *UserDB](db),
	}
}

func (u *userRepo) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[User], error) {
	return u.FindAllPaginatedBy(ctx, pageable, gormen.NewWhere(where.Equal("enable", true)).Build())
}

func (u *userRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
	return u.FindBy(ctx, gormen.NewWhere(where.Equal("username", username)).Build())
}

func main() {
	ctx := context.Background()

	repo := NewUserRepo(SetupDB())

	user1 := User{
		Username: "batch1",
		Password: "1234",
		Enable:   true,
	}
	user2 := User{
		Username: "batch2",
		Password: "1234",
		Enable:   true,
	}
	user3 := User{
		Username: "batch3",
		Password: "1234",
		Enable:   true,
	}
	users := []User{user1, user2, user3}

	err := repo.CreateAll(ctx, &users, 3)
	if err != nil {
		log.Fatalf("Error creating users %v", err)
	}

	user, _ := repo.FindByUsername(ctx, "batch1")
	log.Printf("%+v", user)

	pageRequest, _ := pagination.PageRequestFrom(1, 10, pagination.WithSortOrder("username", sort.Descending))
	page, err := repo.FindAll(ctx, pageRequest)
	if err != nil {
		log.Fatalf("Error searching users %v", err)
	}
	log.Printf("%+v", page)
}

type User struct {
	ID       uint
	Username string
	Password string
	Enable   bool
}

type UserDB struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
	Enable   bool
}

func (udb UserDB) TableName() string {
	return "users"
}

func (udb *UserDB) From(u User) {
	udb.ID = u.ID
	udb.Username = u.Username
	udb.Password = u.Password
	udb.Enable = u.Enable
}

func (udb UserDB) Into() User {
	return User{
		ID:       udb.ID,
		Username: udb.Username,
		Password: udb.Password,
		Enable:   udb.Enable,
	}
}

func SetupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&UserDB{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

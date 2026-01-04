# gormen
*Go library for enhancing Gorm (CRUD, Hexagonal Architecture, Paging, Sorting and Filtering)*

## Description
This library provides a high-level, ergonomic API for enhancing [Gorm](https://github.com/go-gorm/gorm). 
It offers a CRUD API, Paging, Sorting and Filtering. Allowing simple implementation or hexagonal architecture decoupling interface.
Roughly inspired in Java **CrudRepository** and **PagingAndSortingRepository** interfaces.

## Caveats
- This library requires Go 1.23+

## Intallation
```bash
go get -u github.com/javiorfo/gormen@latest
```

## Example Std Repository
#### More examples [here](https://github.com/javiorfo/gormen/tree/master/examples)
```go
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

  user, _ := repo.FindBy(ctx, gormen.NewWhere(where.Like("username", "2%")).Build(), "Person")
  log.Printf("%+v", user)

  orders := []sort.Order{sort.NewOrder("username", sort.Descending)}
  users, err = repo.FindAllOrdered(ctx, orders, "Person")
  if err != nil {
    log.Fatalf("Error searching users %v", err)
  }
  log.Printf("%+v", users)
}
```

## Example Converter Repository (Hexagonal Arch)
#### More examples [here](https://github.com/javiorfo/gormen/tree/master/examples)
```go
// Converter data must satisfies converter interface

// User model
type User struct {
  ID       uint
  Username string
  Password string
  Enable   bool
}

// User DB
type UserDB struct {
  ID       uint   `gorm:"primaryKey;autoIncrement"`
  Username string `gorm:"not null"`
  Password string `gorm:"not null"`
  Enable   bool
}

func (udb UserDB) TableName() string {
  return "users"
}

// Satisfies converter interface
func (udb *UserDB) From(u User) {
  udb.ID = u.ID
  udb.Username = u.Username
  udb.Password = u.Password
  udb.Enable = u.Enable
}

// Satisfies converter interface
func (udb UserDB) Into() User {
  return User{
    ID:       udb.ID,
    Username: udb.Username,
    Password: udb.Password,
    Enable:   udb.Enable,
  }
}

// Then the implementation
type userRepo struct {
  gormen.Repository[User]
}

func NewUserRepo(db *gorm.DB) userRepo {
  return userRepo{
    // converter.Repository converts inside from UserDB to User and vice versa
    // Applying From(M) or Into() M, as needed
    Repository: converter.NewRepository[UserDB, *UserDB](db),
  }
}

func (u *userRepo) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[User], error) {
  return u.FindAllPaginated(ctx, pageable)
}

func (u *userRepo) FindByUsername(ctx context.Context, username string) (*User, error) {
  return u.FindBy(ctx, gormen.NewWhere(where.Equal("username", username)).Build())
}

```

## Available interfaces
#### Any of these satisfies std.Repository or converter.Repository 
```go
// Repository defines a generic interface combining CUD (Create, Update, Delete)
// and Read operations for a model type M.
type Repository[M any] interface {
  CudRepository[M]
  ReadRepository[M]
}

// CudRepository defines generic Create, Update, and Delete operations for model M.
type CudRepository[M any] interface {
  Create(ctx context.Context, model *M) error
  CreateAll(ctx context.Context, model *[]M, batchSize int) error
  Delete(ctx context.Context, model *M) error
  DeleteAll(ctx context.Context, model []M) error
  DeleteAllBy(ctx context.Context, where Where) error
  Save(ctx context.Context, model *M) error
  SaveAll(ctx context.Context, model []M) error
}

// ReadRepository defines generic read/query operations for model M.
type ReadRepository[M any] interface {
  Count(ctx context.Context) (int64, error)
  CountBy(ctx context.Context, where Where) (int64, error)
  FindBy(ctx context.Context, where Where, preloads ...Preload) (*M, error)
  FindAll(ctx context.Context, preloads ...Preload) ([]M, error)
  FindAllBy(ctx context.Context, where Where, preloads ...Preload) ([]M, error)
  FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (*pagination.Page[M], error)
  FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where Where, preloads ...Preload) (*pagination.Page[M], error)
  FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...Preload) ([]M, error)
}
```

---

### Donate
- **Bitcoin** [(QR)](https://raw.githubusercontent.com/javiorfo/img/master/crypto/bitcoin.png)  `1GqdJ63RDPE4eJKujHi166FAyigvHu5R7v`
- [Paypal](https://www.paypal.com/donate/?hosted_button_id=FA7SGLSCT2H8G)

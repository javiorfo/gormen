package pagination

import (
	"github.com/javiorfo/gormen/pagination/sort"
	"gorm.io/gorm"
)

type Page[T any] struct {
	Total    int64
	Elements []T
}

type Pageable interface {
	PageNumber() int
	PageSize() int
	SortOrders() []sort.Order
	Paginate(*gorm.DB) (*gorm.DB, error)
	Order(*gorm.DB) (*gorm.DB, error)
	Filter(*gorm.DB) (*gorm.DB, error)
}

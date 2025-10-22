package pagination

import (
	"github.com/javiorfo/gormen/pagination/sort"
	"gorm.io/gorm"
)

// Page represents a paginated result containing total item count and current page elements.
type Page[T any] struct {
	// Total number of items available (no paging involved)
	Total int64
	// Current page elements
	Elements []T
}

// Pageable defines an interface for pagination, sorting, and filtering capabilities
// that can be applied to a GORM DB query.
type Pageable interface {
	// PageNumber returns the current page number (starting from 1).
	PageNumber() int
	// PageSize returns the number of items per page.
	PageSize() int
	// SortOrders returns a list of sort.Order specifying sorting criteria.
	SortOrders() []sort.Order
	// Paginate applies pagination limits and offsets to the GORM DB instance.
	Paginate(*gorm.DB) (*gorm.DB, error)
	// Order applies sorting orders to the GORM DB instance.
	Order(*gorm.DB) (*gorm.DB, error)
	// Filter applies filtering criteria to the GORM DB instance.
	Filter(*gorm.DB) (*gorm.DB, error)
}

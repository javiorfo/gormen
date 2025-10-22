package pagination

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

// pageRequest holds pagination, sorting, and filtering data for database queries.
type pageRequest struct {
	// The current page number, starting from 1.
	pageNumber int
	// The number of items per page.
	pageSize int
	// List of sorting criteria for query results.
	sortOrders []sort.Order
	// Optional filter criteria to narrow results.
	filter nilo.Option[any]
}

// PageNumber returns the current page number.
func (p *pageRequest) PageNumber() int {
	return p.pageNumber
}

// PageSize returns the number of items per page.
func (p *pageRequest) PageSize() int {
	return p.pageSize
}

// SortOrders returns the list of sorting criteria.
func (p *pageRequest) SortOrders() []sort.Order {
	return p.sortOrders
}

// Paginate applies offset and limit based on page number and size to the GORM DB query.
func (p *pageRequest) Paginate(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(p.order(p.paginate(db)), p.filter)
}

// Order applies sorting based on the sortOrders field to the GORM DB query.
func (p *pageRequest) Order(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(p.order(db), p.filter)
}

// Filter applies filtering criteria to the GORM DB query.
func (p *pageRequest) Filter(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(db, p.filter)
}

// paginate modifies the given GORM DB query with offset and limit for pagination.
func (p *pageRequest) paginate(db *gorm.DB) *gorm.DB {
	return db.Offset((p.pageNumber - 1) * p.pageSize).Limit(p.pageSize)
}

// order applies sorting orders to the GORM DB query.
func (p *pageRequest) order(db *gorm.DB) *gorm.DB {
	for _, o := range p.sortOrders {
		db = db.Order(o.Get())
	}
	return db
}

// DefaultPageRequest returns a pageRequest with default settings.
func DefaultPageRequest() *pageRequest {
	return &pageRequest{
		pageNumber: 1,
		pageSize:   10,
		sortOrders: []sort.Order{sort.Default()},
		filter:     nilo.None[any](),
	}
}

// PageOptions is a function that modifies a pageRequest, often used for optional configurations.
type PageOptions func(*pageRequest) error

// WithSortOrder adds a sorting order to the pageRequest.
func WithSortOrder(by string, direction sort.Direction) PageOptions {
	order := sort.NewOrder(by, direction)
	return func(p *pageRequest) error {
		if err := order.IsValid(); err != nil {
			return err
		}
		p.sortOrders = append(p.sortOrders, order)
		return nil
	}
}

// WithFilter adds a filter struct to the pageRequest.
func WithFilter(filter any) PageOptions {
	return func(p *pageRequest) error {
		if reflect.TypeOf(filter).Kind() != reflect.Struct {
			return errors.New("'filter' must be a struct")
		}
		p.filter = nilo.Some(filter)
		return nil
	}
}

// PageRequestFrom constructs a pageRequest from given page number, page size, and options.
func PageRequestFrom[T interface{ ~int | ~string }](pageNumber, pageSize T, options ...PageOptions) (*pageRequest, error) {
	pageNumberInt, err := toInt(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, errors.New("'pageNumber' must be a positive number")
	}

	pageSizeInt, err := toInt(pageSize)
	if err != nil || pageSizeInt < 1 {
		return nil, errors.New("'pageSize' must be a positive number greater than 0")
	}

	if pageSizeInt < pageNumberInt {
		return nil, errors.New("'pageSize' must be greater than 'pageNumber'")
	}

	p := &pageRequest{}

	for _, opt := range options {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}

	p.pageNumber = pageNumberInt
	p.pageSize = pageSizeInt

	return p, nil
}

// toInt converts a value of type int or string to int.
func toInt[T interface{ ~int | ~string }](value T) (int, error) {
	switch v := any(value).(type) {
	case int:
		return v, nil
	case string:
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	default:
		return 0, errors.New("not supported")
	}
}

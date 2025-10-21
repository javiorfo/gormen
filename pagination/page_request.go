package pagination

import (
	"errors"
	"reflect"

	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

type pageRequest struct {
	pageNumber int
	pageSize   int
	sortOrders []sort.Order
	filter     nilo.Option[any]
}

func (p pageRequest) PageNumber() int {
	return p.pageNumber
}

func (p pageRequest) PageSize() int {
	return p.pageSize
}

func (p pageRequest) SortOrders() []sort.Order {
	return p.sortOrders
}

func (p pageRequest) Paginate(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(p.order(p.paginate(db)), p.filter)
}

func (p pageRequest) Order(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(p.order(db), p.filter)
}

func (p pageRequest) Filter(db *gorm.DB) (*gorm.DB, error) {
	return filterValues(db, p.filter)
}

func (p pageRequest) paginate(db *gorm.DB) *gorm.DB {
	db = db.Offset(p.pageNumber - 1).Limit(p.pageSize)
	return db
}

func (p pageRequest) order(db *gorm.DB) *gorm.DB {
	for _, o := range p.sortOrders {
		db = db.Order(o.Prepared())
	}
	return db
}

func DefaultPageRequest() pageRequest {
	return pageRequest{
		pageNumber: 1,
		pageSize:   10,
		sortOrders: []sort.Order{sort.Default()},
		filter:     nilo.None[any](),
	}
}

type PageOptions func(*pageRequest) error

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

func WithFilter(filter any) PageOptions {
	return func(p *pageRequest) error {
		if reflect.TypeOf(filter).Kind() != reflect.Struct {
			return errors.New("'filter' must be a struct")
		}

		p.filter = nilo.Some(filter)
		return nil
	}
}

func PageRequestFrom[T interface{ ~int | ~string }](pageNumber, pageSize T, options ...PageOptions) (*pageRequest, error) {
	pageNumberInt, ok := any(pageNumber).(int)
	if !ok || pageNumberInt < 0 {
		return nil, errors.New("'pageNumber' must be a positive number")
	}

	pageSizeInt, ok := any(pageSize).(int)
	if !ok || pageSizeInt < 1 {
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

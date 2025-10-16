package pagination

import (
	"errors"

	"github.com/javiorfo/gormix/pagination/sort"
	"github.com/javiorfo/gormix/response"
)

type Page struct {
	Page       int
	Size       int
	SortOrders []sort.Order
}

func DefaultPage() Page {
	return Page{
		Page:       1,
		Size:       10,
		SortOrders: []sort.Order{sort.NewOrder("id", sort.Ascending)},
	}
}

func NewPage[T any](page, size T, sortOrders ...sort.Order) (*Page, error) {
	pageInt, ok := any(page).(int)
	if !ok || pageInt < 0 {
		return nil, errors.New("'page' must be a positive number")
	}

	sizeInt, ok := any(size).(int)
	if !ok || sizeInt < 0 {
		return nil, errors.New("'size' must be a positive number")
	}

	if sizeInt < pageInt {
		return nil, errors.New("'size' must be greater than 'page'")
	}

	for _, order := range sortOrders {
		if err := order.IsValid(); err != nil {
			return nil, err
		}
	}

	return &Page{pageInt, sizeInt, sortOrders}, nil
}

func Paginator(p Page, total int64) response.Pagination {
	return response.Pagination{
		PageNumber: p.Page,
		PageSize:   p.Size,
		Total:      total,
	}
}

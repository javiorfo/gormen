package sort

import "errors"

type Order struct {
	by        string
	direction Direction
}

func (o Order) By() string {
	return o.by
}

func (o Order) Direction() Direction {
	return o.direction
}

func (o Order) IsValid() error {
	if o.direction != Ascending && o.direction != Descending {
		return errors.New("'order.direction' must be 'asc' or 'desc'")
	}

	if o.by == "" {
		return errors.New("'order.by' must not be empty")
	}

	return nil
}

func NewOrder(by string, direction Direction) Order {
	return Order{by, direction}
}

type Direction = string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

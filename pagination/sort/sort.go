package sort

import (
	"errors"
	"fmt"
	"strings"
)

// Order represents sorting criteria with a column and direction.
type Order struct {
	// Column name to sort by
	by string
	// Sort direction, either Ascending or Descending
	direction Direction
}

// By returns the column name used for ordering.
func (o Order) By() string {
	return o.by
}

// Direction returns the sorting direction.
func (o Order) Direction() Direction {
	return o.direction
}

// Get returns the SQL fragment for the order clause (e.g., "column asc").
func (o Order) Get() string {
	return fmt.Sprintf("%s %s", o.By(), o.Direction())
}

// IsValid validates the Order, ensuring direction is 'asc' or 'desc' and column is not empty.
func (o Order) IsValid() error {
	if o.direction != Ascending && o.direction != Descending {
		return errors.New("'order.direction' must be 'asc' or 'desc'")
	}

	if o.by == "" {
		return errors.New("'order.by' must not be empty")
	}

	return nil
}

// Default returns a default Order by "id" in ascending order.
func Default() Order {
	return NewOrder("id", Ascending)
}

// NewOrder creates a new Order with the given column and direction.
func NewOrder(by string, direction Direction) Order {
	return Order{by, direction}
}

// Direction defines the sorting direction as a string type.
type Direction = string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

// DirectionFromString converts a string to a Direction, defaulting to Ascending.
func DirectionFromString(dir string) Direction {
	if strings.EqualFold(dir, Descending) {
		return Descending
	}
	return Ascending
}

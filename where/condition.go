package where

import (
	"fmt"

	"github.com/javiorfo/gormen/internal/utils"
)

// ColumnName represents the name of a database column.
type ColumnName = string

// Value represents a value to match in a condition; it can be any type.
type Value = any

// Condition defines an interface for query conditions,
// requiring a Get method that returns a query string and value.
type Condition interface {
	Get() (string, any)
}

// like represents a SQL LIKE condition for pattern matching.
type like struct {
	name  ColumnName
	value Value
}

// Like constructs a LIKE condition for the given column name and value.
func Like(name ColumnName, value Value) like {
	return like{name, value}
}

// Get returns the SQL LIKE query snippet and its corresponding value.
// Satisfies Condition interface
func (l like) Get() (string, any) {
	return fmt.Sprintf("%s like ?", l.name), l.value
}

// in represents a SQL IN condition to match a column against multiple values.
type in struct {
	name  ColumnName
	value Value
}

// In constructs an IN condition for the given column name and values.
func In(name ColumnName, value Value) in {
	return in{name, value}
}

// Get returns the SQL IN query snippet and the processed list of values.
// Satisfies Condition interface
func (i in) Get() (string, any) {
	v := i.value
	utils.GetValueAsCommaSeparated(i.value).Consume(func(s []string) {
		v = s
	})

	return fmt.Sprintf("%s in (?)", i.name), v
}

// equal represents a SQL equality condition.
type equal struct {
	name  ColumnName
	value Value
}

// Equal constructs an equality condition for the given column name and value.
func Equal(name ColumnName, value Value) equal {
	return equal{name, value}
}

// Get returns the SQL equality query snippet and its corresponding value.
// Satisfies Condition interface
func (e equal) Get() (string, any) {
	return fmt.Sprintf("%s = ?", e.name), e.value
}

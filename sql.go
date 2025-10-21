package gormen

import (
	"fmt"

	"github.com/javiorfo/gormen/internal"
)

type Preload = string

type ColumnName string
type Value any

type sqlField struct {
	name  ColumnName
	value Value
}

func (s sqlField) Prepared() string {
	return fmt.Sprintf("%s = ?", s.name)
}

func (s sqlField) Name() ColumnName {
	return s.name
}

func (s sqlField) Value() Value {
	return s.value
}

type Where map[sqlField]int

func NewWhere(name ColumnName, value Value) Where {
	sf := sqlField{name, value}
	return map[sqlField]int{sf: internal.None}
}

func (w Where) And(name ColumnName, value Value) {
	sf := sqlField{name, value}
	w[sf] = internal.And
}

func (w Where) Or(name ColumnName, value Value) {
	sf := sqlField{name, value}
	w[sf] = internal.Or
}

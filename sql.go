package gormen

import (
	"fmt"

	"github.com/javiorfo/gormen/internal/types"
)

type Preload = string
type Join = string
type Conditions = map[sqlField]int

type ColumnName = string
type Value = any

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

type Where struct {
	conditions Conditions
	joins      []Join
}

func (w Where) Conditions() Conditions {
	return w.conditions
}

func (w Where) Joins() []Join {
	return w.joins
}

func NewWhere(name ColumnName, value Value) *Where {
	sf := sqlField{name, value}
	return &Where{conditions: Conditions{sf: types.None}}
}

func (w *Where) And(name ColumnName, value Value) *Where {
	sf := sqlField{name, value}
	w.conditions[sf] = types.And
	return w
}

func (w *Where) Or(name ColumnName, value Value) *Where {
	sf := sqlField{name, value}
	w.conditions[sf] = types.Or
	return w
}

func (w *Where) WithJoin(joins ...string) *Where {
	w.joins = joins
	return w
}

func (w *Where) Build() Where {
	return *w
}

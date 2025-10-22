package gormen

import (
	"fmt"

	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/gormen/internal/utils"
)

type Preload = string
type Join = string
type Conditions = map[Condition]int

type ColumnName = string
type Value = any

type equal struct {
	name  ColumnName
	value Value
}

func Equal(name ColumnName, value Value) equal {
	return equal{name, value}
}

type like struct {
	name  ColumnName
	value Value
}

func (l like) Get() (string, any) {
	return fmt.Sprintf("%s like ?", l.name), l.value
}

type in struct {
	name  ColumnName
	value Value
}

func In(name ColumnName, value Value) in {
	return in{name, value}
}

func (i in) Get() (string, any) {
	value := utils.GetValueAsCommaSeparated(i.value).
		MapOrAny(i.value, func(s []string) any {
			return s
		})

	return fmt.Sprintf("%s in ?", i.name), value
}

type Condition interface {
	Get() (string, any)
}

func (e equal) Get() (string, any) {
	return fmt.Sprintf("%s = ?", e.name), e.value
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

func NewWhere(c Condition) *Where {
	return &Where{conditions: Conditions{c: types.None}}
}

func (w *Where) And(c Condition) *Where {
	w.conditions[c] = types.And
	return w
}

func (w *Where) Or(c Condition) *Where {
	w.conditions[c] = types.Or
	return w
}

func (w *Where) WithJoin(joins ...string) *Where {
	w.joins = joins
	return w
}

func (w *Where) Build() Where {
	return *w
}

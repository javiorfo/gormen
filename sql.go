package gormen

import "fmt"

type Preload = string

type ColumnName string
type Value any

type SqlField struct {
	name  ColumnName
	value Value
}

func (s SqlField) Prepared() string {
	return fmt.Sprintf("%s = ?", s.name)
}

func (s SqlField) Name() ColumnName {
	return s.name
}

func (s SqlField) Value() Value {
	return s.value
}

func NewSqlField(name ColumnName, value Value) SqlField {
	return SqlField{name, value}
}

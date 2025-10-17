package crud

import (
	"context"
	"fmt"

	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

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

func NewSqlField(name ColumnName, value Value) SqlField {
	return SqlField{name, value}
}

type hexaRepository[E any, C converter[E, M], M any] struct {
	db *gorm.DB
}

type converter[E, M any] interface {
	*E
	From(M)
	Into() M
}

func NewHexaRepository[E any, C converter[E, M], M any](db *gorm.DB) *hexaRepository[E, C, M] {
	return &hexaRepository[E, C, M]{db}
}

type stdRepository[M any] struct {
	db *gorm.DB
}

type Repository[M any] interface {
	FindBy(ctx context.Context, sqlField SqlField, preloads ...Preload) (nilo.Option[M], error)
}

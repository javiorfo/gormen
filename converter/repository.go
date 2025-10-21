package converter

import (
	"github.com/javiorfo/gormen"
	"gorm.io/gorm"
)

type repository[E any, C converter[E, M], M any] struct {
	db *gorm.DB
}

type converter[E, M any] interface {
	*E
	From(M)
	Into() M
}

func NewRepository[E any, C converter[E, M], M any](db *gorm.DB) gormen.Repository[M] {
	return &repository[E, C, M]{db}
}

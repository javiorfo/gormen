package converter

import "gorm.io/gorm"

type converterRepository[E any, C converter[E, M], M any] struct {
	db *gorm.DB
}

type converter[E, M any] interface {
	*E
	From(M)
	Into() M
}

func NewConverterRepository[E any, C converter[E, M], M any](db *gorm.DB) *converterRepository[E, C, M] {
	return &converterRepository[E, C, M]{db}
}

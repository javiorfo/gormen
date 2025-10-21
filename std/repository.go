package std

import (
	"github.com/javiorfo/gormen"
	"gorm.io/gorm"
)

type repository[M any] struct {
	db *gorm.DB
}

func NewRepository[M any](db *gorm.DB) gormen.Repository[M] {
	return &repository[M]{db}
}

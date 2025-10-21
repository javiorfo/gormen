package std

import "gorm.io/gorm"

type repository[M any] struct {
	db *gorm.DB
}

func NewRepository[M any](db *gorm.DB) *repository[M] {
	return &repository[M]{db}
}

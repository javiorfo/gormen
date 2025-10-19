package std

import "gorm.io/gorm"

type stdRepository[M any] struct {
	db *gorm.DB
}

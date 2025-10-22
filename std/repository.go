package std

import (
	"github.com/javiorfo/gormen"
	"gorm.io/gorm"
)

// repository is a generic struct implementing the Repository interface using GORM.
type repository[M any] struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance for model M with the given GORM DB.
func NewRepository[M any](db *gorm.DB) gormen.Repository[M] {
	return &repository[M]{db}
}

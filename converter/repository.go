package converter

import (
	"github.com/javiorfo/gormen"
	"gorm.io/gorm"
)

// repository is a generic struct implementing Repository with entity E,
// a converter C to transform between entity E and model M, using GORM DB.
type repository[E any, C converter[E, M], M any] struct {
	db *gorm.DB // GORM database connection
}

// converter defines an interface for types that can convert between
// entity E (pointer receiver) and model M.
type converter[E, M any] interface {
	// Pointer to DB entity type E as receiver
	*E
	// Populate DB entity from model M
	From(M)
	// Convert DB entity into model M
	Into() M
}

// NewRepository creates a new repository for entity E, converter C, and model M using GORM DB.
func NewRepository[E any, C converter[E, M], M any](db *gorm.DB) gormen.Repository[M] {
	return &repository[E, C, M]{db}
}

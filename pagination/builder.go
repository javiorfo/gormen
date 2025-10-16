package pagination

import (
	"fmt"

	"gorm.io/gorm"
)

func paginate(db *gorm.DB, p Page) *gorm.DB {
	db = db.Offset(p.Page - 1).Limit(p.Size)

	for _, o := range p.SortOrders {
		db = db.Order(fmt.Sprintf("%s %s", o.By(), o.Direction()))
	}
	return db
}

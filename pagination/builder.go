package pagination

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

type tagAndValue struct {
	tagValue   string
	fieldValue any
	joins      string
}

func paginate(db *gorm.DB, p pageRequest) *gorm.DB {
	db = db.Offset(p.pageNumber - 1).Limit(p.pageSize)

	for _, o := range p.sortOrders {
		db = db.Order(fmt.Sprintf("%s %s", o.By(), o.Direction()))
	}
	return db
}

func filterValues(db *gorm.DB, filter nilo.Option[any]) (*gorm.DB, error) {
	if filter.IsNone() {
		return db, nil
	}

	var results []tagAndValue

	v := reflect.ValueOf(filter.Unwrap())
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		tagValue := field.Tag.Get("filter")
		if fmt.Sprintf("%v", value.Interface()) == "" || tagValue == "" {
			continue
		}

		parts := strings.Split(tagValue, ";")
		// TODO validate array
		filterString := parts[0]

		var joins string
		for _, part := range parts[1:] {
			if after, ok := strings.CutPrefix(part, "join:"); ok {
				joins = fmt.Sprintf("%s %s ", joins, after)
			}
		}

		results = append(results, tagAndValue{
			tagValue:   filterString,
			fieldValue: value.Interface(),
			joins:      joins,
		})
	}

	for _, v := range results {
		if v.joins != "" {
			db = db.Joins(v.joins)
		}
		db = db.Where(v.tagValue, v.fieldValue)
	}

	return db, nil
}

package pagination

import (
	"errors"
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
		if tagValue == "" {
			return db, errors.New("'filter' tag must exist in all properties")
		}

		if fmt.Sprintf("%v", value.Interface()) == "" {
			continue
		}

		parts := strings.Split(tagValue, ";")
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

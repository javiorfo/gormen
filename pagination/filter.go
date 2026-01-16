package pagination

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/javiorfo/gormen/internal/utils"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

// tagAndValue holds a parsed filter tag, its corresponding field value,
// and any SQL join clauses associated with the filter.
type tagAndValue struct {
	tagValue   string
	fieldValue any
	joins      string
}

// filterValues applies filtering conditions from a struct with "filter" tags to a GORM DB query.
// It takes an optional filter struct wrapped in nilo.Option and returns the modified DB instance.
// Returns an error if any struct field lacks the "filter" tag.
func filterValues(db *gorm.DB, filter nilo.Option[any]) (*gorm.DB, error) {
	if filter.IsNil() {
		return db, nil
	}

	var results []tagAndValue

	v := reflect.ValueOf(filter.AsValue())
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

		fieldValue := value.Interface()
		utils.GetValueAsCommaSeparated(fieldValue).Consume(func(s []string) {
			fieldValue = s
		})
		results = append(results, tagAndValue{
			tagValue:   filterString,
			fieldValue: fieldValue,
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

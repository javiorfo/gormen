package pagination

import (
	"errors"
	"testing"

	"github.com/javiorfo/gormen/pagination/sort"
)

func TestWithSortOrder_ValidAndInvalid(t *testing.T) {
	p := &pageRequest{}
	err := WithSortOrder("name", sort.Ascending)(p)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(p.sortOrders) != 1 {
		t.Errorf("expected 1 sort order, got %d", len(p.sortOrders))
	}

	err = WithSortOrder("name", "garbage")(p)
	if err == nil {
		t.Error("expected error for invalid sort order direction")
	}
}

func TestWithFilter_ValidAndInvalid(t *testing.T) {
	type filterStruct struct {
		Field string
	}

	p := &pageRequest{}

	err := WithFilter(filterStruct{Field: "test"})(p)
	if err != nil {
		t.Errorf("unexpected error for valid filter struct: %v", err)
	}
	if p.filter.IsNil() {
		t.Error("expected filter option to be set")
	}

	err = WithFilter("not a struct")(p)
	if err == nil {
		t.Error("expected error when passing non-struct filter")
	}
}

func TestPageRequestFrom_ValidAndInvalid(t *testing.T) {
	p, err := PageRequestFrom("1", "10")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.pageNumber != 1 || p.pageSize != 10 {
		t.Errorf("unexpected pageNumber or pageSize")
	}

	_, err = PageRequestFrom(-1, 10)
	if err == nil {
		t.Error("expected error for negative pageNumber")
	}

	_, err = PageRequestFrom(1, 0)
	if err == nil {
		t.Error("expected error for pageSize zero")
	}

	_, err = PageRequestFrom(5, 3)
	if err == nil {
		t.Error("expected error for pageSize < pageNumber")
	}
}

func TestPageRequestFrom_WithOptions(t *testing.T) {
	optApplied := false
	opt := func(p *pageRequest) error {
		optApplied = true
		return nil
	}
	p, err := PageRequestFrom(1, 10, opt)
	if err != nil {
		t.Errorf("unexpected error applying options: %v", err)
	}
	if !optApplied {
		t.Error("expected option function to be applied")
	}
	if p.pageNumber != 1 || p.pageSize != 10 {
		t.Errorf("unexpected pageNumber or pageSize")
	}
}

func TestPageRequestFrom_OptionReturnsError(t *testing.T) {
	opt := func(p *pageRequest) error {
		return errors.New("option error")
	}
	_, err := PageRequestFrom(1, 10, opt)
	if err == nil {
		t.Error("expected error from option")
	}
	if err.Error() != "option error" {
		t.Errorf("expected 'option error', got %v", err)
	}
}

package pagination

import (
	"testing"

	"github.com/javiorfo/gormen/internal/testutils"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	db = testutils.SetupTestDB()
	m.Run()
}

type testFilter struct {
	Name  string `filter:"name = ?"`
	IDs   []int  `filter:"id IN ?;join:inner join users u on u.id = profiles.user_id"`
	Empty string `filter:"empty = ?"`
}

func TestFilterValues_Nil(t *testing.T) {
	filterOpt := nilo.Nil[any]()
	gotDB, err := filterValues(db, filterOpt)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if gotDB != db {
		t.Errorf("expected original db instance when filter is none")
	}
}

func TestFilterValues_EmptyFields(t *testing.T) {
	filterOpt := nilo.Value(any(testFilter{Name: "", IDs: nil}))
	gotDB, err := filterValues(db, filterOpt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotDB == nil {
		t.Fatalf("expected non-nil Gorm DB")
	}
}

func TestFilterValues_FilterWithJoins(t *testing.T) {
	filterOpt := nilo.Value(any(testFilter{Name: "alice", IDs: []int{1, 2, 3}}))
	gotDB, err := filterValues(db, filterOpt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotDB == nil {
		t.Fatalf("expected non-nil Gorm DB")
	}
}

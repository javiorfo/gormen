package converter

import (
	"context"
	"testing"

	"github.com/javiorfo/gormen"
)

func TestFindBy(t *testing.T) {
	ctx := context.Background()

	optional, err := repo.FindBy(ctx, gormen.NewWhere("username", "jdoe"))
	if err != nil {
		t.Fatal("finding by username")
	}

	if optional.IsNone() {
		t.Fatal("user not found")
	}
}

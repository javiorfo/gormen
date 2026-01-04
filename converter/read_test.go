package converter

import (
	"context"
	"testing"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/gormen/where"
)

func TestRead(t *testing.T) {
	ctx := context.Background()

	t.Run("Converter FindBy", func(t *testing.T) {
		value, err := repo.FindBy(ctx, gormen.NewWhere(where.Equal("persons.name", "Batch 2")).
			WithJoin("inner join persons on users.person_id = persons.id").Build(), "Person")

		if err != nil {
			t.Fatalf("executing find by username %v\n", err)
		}

		if value == nil {
			t.Fatal("user not found")
		}

		if value.Password != "123" {
			t.Fatal("user found does not match")
		}
	})

	t.Run("Converter FindBy Not Found", func(t *testing.T) {
		value, err := repo.FindBy(ctx, gormen.NewWhere(where.Equal("username", "notfound")).Build())
		if err != nil {
			t.Fatalf("executing find by username %v\n", err)
		}

		if value != nil {
			t.Fatal("user must be nil")
		}
	})

	t.Run("Converter FindAll", func(t *testing.T) {
		users, err := repo.FindAll(ctx)

		if err != nil {
			t.Fatalf("executing find all %v\n", err)
		}

		if len(users) != 3 {
			t.Fatalf("len must be 3, got %d\n", len(users))
		}
	})

	t.Run("Converter FindAllBy", func(t *testing.T) {
		users, err := repo.FindAllBy(ctx, gormen.NewWhere(where.In("password", "123,1234")).Build())

		if err != nil {
			t.Fatalf("executing find all by %v\n", err)
		}

		if len(users) != 3 {
			t.Fatalf("len must be 3, got %d\n", len(users))
		}
	})

	t.Run("Converter Count", func(t *testing.T) {
		count, err := repo.Count(ctx)

		if err != nil {
			t.Fatalf("executing count %v\n", err)
		}

		if count != 3 {
			t.Fatalf("count must be 3, got %d\n", count)
		}
	})

	t.Run("Converter CountBy", func(t *testing.T) {
		count, err := repo.CountBy(ctx, gormen.NewWhere(where.Equal("password", "123")).Build())

		if err != nil {
			t.Fatalf("executing count by %v\n", err)
		}

		if count != 2 {
			t.Fatalf("count must be 2, got %d\n", count)
		}
	})

	t.Run("Converter FindAllPaginated default", func(t *testing.T) {
		pageRequest := pagination.DefaultPageRequest()
		page, err := repo.FindAllPaginated(ctx, pageRequest, "Person")

		if err != nil {
			t.Fatalf("executing find all paginated default %v\n", err)
		}

		if page.Total != 3 {
			t.Fatalf("total must be 3, got %d\n", page.Total)
		}
	})

	t.Run("Converter FindAllPaginated with page", func(t *testing.T) {
		pageRequest, err := pagination.PageRequestFrom(1, 2)
		if err != nil {
			t.Fatalf("creating page request %v\n", err)
		}

		page, err := repo.FindAllPaginated(ctx, pageRequest, "Person")
		if err != nil {
			t.Fatalf("executing find all paginated with page %v\n", err)
		}

		if page.Total != 3 {
			t.Fatalf("total must be 3, got %d\n", page.Total)
		}

		if len(page.Elements) != 2 {
			t.Fatalf("len must be 2, got %d\n", len(page.Elements))
		}
	})

	t.Run("Converter FindAllPaginated with page, sort and filter", func(t *testing.T) {
		type UserFilter struct {
			Ids string `filter:"persons.id in (?);join:inner join persons on users.person_id = persons.id"`
		}

		pageRequest, err := pagination.PageRequestFrom(1, 1,
			pagination.WithSortOrder("username", sort.DirectionFromString("DESC")),
			pagination.WithFilter(UserFilter{"1,2"}),
		)

		if err != nil {
			t.Fatalf("creating page request %v\n", err)
		}

		page, err := repo.FindAllPaginated(ctx, pageRequest, "Person")
		if err != nil {
			t.Fatalf("executing find all paginated %v\n", err)
		}

		if page.Total != 2 {
			t.Fatalf("total must be 2, got %d\n", page.Total)
		}

		if len(page.Elements) != 1 {
			t.Fatalf("len must be 1, got %d\n", len(page.Elements))
		}

		if page.Elements[0].Username != "jdoe" {
			t.Fatalf("sorting elements. Got %s\n", page.Elements[0].Username)
		}
	})
}

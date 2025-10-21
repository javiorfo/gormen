package converter

import (
	"context"
	"testing"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
)

func TestRead(t *testing.T) {
	ctx := context.Background()

	t.Run("FindBy", func(t *testing.T) {
		optional, err := repo.FindBy(ctx, gormen.NewWhere("persons.name", "Batch 2").
			WithJoin("inner join persons on users.person_id = persons.id").Build(), "Person")

		if err != nil {
			t.Fatalf("executing find by username %v\n", err)
		}

		if optional.IsNone() {
			t.Fatal("user not found")
		}

		if optional.Unwrap().Password != "123" {
			t.Fatal("user found does not match")
		}
	})

	t.Run("FindBy Not Found", func(t *testing.T) {
		optional, err := repo.FindBy(ctx, gormen.NewWhere("username", "notfound").Build())
		if err != nil {
			t.Fatalf("executing find by username %v\n", err)
		}

		if optional.IsSome() {
			t.Fatal("user must be None")
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		users, err := repo.FindAll(ctx)

		if err != nil {
			t.Fatalf("executing find all %v\n", err)
		}

		if len(users) != 3 {
			t.Fatalf("len must be 3, got %d\n", len(users))
		}
	})

	t.Run("FindAllBy", func(t *testing.T) {
		users, err := repo.FindAllBy(ctx, gormen.NewWhere("password", "123").Build())

		if err != nil {
			t.Fatalf("executing find all by %v\n", err)
		}

		if len(users) != 2 {
			t.Fatalf("len must be 2, got %d\n", len(users))
		}
	})

	t.Run("Count", func(t *testing.T) {
		count, err := repo.Count(ctx)

		if err != nil {
			t.Fatalf("executing count %v\n", err)
		}

		if count != 3 {
			t.Fatalf("count must be 3, got %d\n", count)
		}
	})

	t.Run("CountBy", func(t *testing.T) {
		count, err := repo.CountBy(ctx, gormen.NewWhere("password", "123").Build())

		if err != nil {
			t.Fatalf("executing count by %v\n", err)
		}

		if count != 2 {
			t.Fatalf("count must be 2, got %d\n", count)
		}
	})

	t.Run("FindAllPaginated default", func(t *testing.T) {
		pageRequest := pagination.DefaultPageRequest()
		page, err := repo.FindAllPaginated(ctx, pageRequest, "Person")

		if err != nil {
			t.Fatalf("executing find all paginated default %v\n", err)
		}

		if page.Total != 3 {
			t.Fatalf("total must be 3, got %d\n", page.Total)
		}
	})

	t.Run("FindAllPaginated with page", func(t *testing.T) {
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

	t.Run("FindAllPaginated with page, sort and filter", func(t *testing.T) {
		type UserFilter struct {
			Name string `filter:"persons.name like ?;join:inner join persons on users.person_id = persons.id"`
		}

		pageRequest, err := pagination.PageRequestFrom(1, 1,
			pagination.WithSortOrder("username", sort.DirectionFromString("DESC")),
			pagination.WithFilter(UserFilter{"%Batch%"}),
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

		if page.Elements[0].Username != "batch2" {
			t.Fatalf("sorting elements. Got %s\n", page.Elements[0].Username)
		}
	})
}

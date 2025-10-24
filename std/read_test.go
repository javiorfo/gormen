package std

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

	t.Run("Std FindBy", func(t *testing.T) {
		optional, err := repo.FindBy(ctx, gormen.NewWhere(where.Equal("persons.name", "Batch 2")).
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

	t.Run("Std FindBy Not Found", func(t *testing.T) {
		optional, err := repo.FindBy(ctx, gormen.NewWhere(where.Equal("username", "notfound")).Build())
		if err != nil {
			t.Fatalf("executing find by username %v\n", err)
		}

		if optional.IsSome() {
			t.Fatal("user must be None")
		}
	})

	t.Run("Std FindAll", func(t *testing.T) {
		users, err := repo.FindAll(ctx)

		if err != nil {
			t.Fatalf("executing find all %v\n", err)
		}

		if len(users) != 3 {
			t.Fatalf("len must be 3, got %d\n", len(users))
		}
	})

	t.Run("Std FindAllBy", func(t *testing.T) {
		users, err := repo.FindAllBy(ctx, gormen.NewWhere(where.Like("username", "batch%")).And(where.In("id", "1,2,3")).Build())

		if err != nil {
			t.Fatalf("executing find all by %v\n", err)
		}

		if len(users) != 2 {
			t.Fatalf("len must be 2, got %d\n", len(users))
		}
	})

	t.Run("Std Count", func(t *testing.T) {
		count, err := repo.Count(ctx)

		if err != nil {
			t.Fatalf("executing count %v\n", err)
		}

		if count != 3 {
			t.Fatalf("count must be 3, got %d\n", count)
		}
	})

	t.Run("Std CountBy", func(t *testing.T) {
		count, err := repo.CountBy(ctx, gormen.NewWhere(where.Equal("password", "123")).Build())

		if err != nil {
			t.Fatalf("executing count by %v\n", err)
		}

		if count != 2 {
			t.Fatalf("count must be 2, got %d\n", count)
		}
	})

	t.Run("Std FindAllPaginated default", func(t *testing.T) {
		pageRequest := pagination.DefaultPageRequest()
		page, err := repo.FindAllPaginated(ctx, pageRequest, "Person")

		if err != nil {
			t.Fatalf("executing find all paginated default %v\n", err)
		}

		if page.Total != 3 {
			t.Fatalf("total must be 3, got %d\n", page.Total)
		}
	})

	t.Run("Std FindAllPaginated with page", func(t *testing.T) {
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

	t.Run("Std FindAllPaginated with page, sort and filter", func(t *testing.T) {
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

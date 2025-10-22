package gormen

import (
	"context"

	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
)

// Repository defines a generic interface combining CUD (Create, Update, Delete)
// and Read operations for a model type M.
type Repository[M any] interface {
	CudRepository[M]
	ReadRepository[M]
}

// CudRepository defines generic Create, Update, and Delete operations for model M.
type CudRepository[M any] interface {
	// Create inserts a new record for the given model.
	Create(ctx context.Context, model *M) error
	// CreateAll inserts multiple records in batches of specified size.
	CreateAll(ctx context.Context, model *[]M, batchSize int) error
	// Delete removes the specified model record.
	Delete(ctx context.Context, model *M) error
	// DeleteAll removes all specified model records.
	DeleteAll(ctx context.Context, model []M) error
	// DeleteAllBy removes all records matching the given Where condition.
	DeleteAllBy(ctx context.Context, where Where) error
	// Save creates or updates the given model record.
	Save(ctx context.Context, model *M) error
	// SaveAll creates or updates multiple model records.
	SaveAll(ctx context.Context, model []M) error
}

// ReadRepository defines generic read/query operations for model M.
type ReadRepository[M any] interface {
	// Count returns the total number of records.
	Count(ctx context.Context) (int64, error)
	// CountBy returns the number of records matching a specific condition.
	CountBy(ctx context.Context, where Where) (int64, error)
	// FindBy returns a single record matching the condition or none if not found.
	FindBy(ctx context.Context, where Where, preloads ...Preload) (nilo.Option[M], error)
	// FindAll returns all records.
	FindAll(ctx context.Context, preloads ...Preload) ([]M, error)
	// FindAllBy returns all records matching a condition.
	FindAllBy(ctx context.Context, where Where, preloads ...Preload) ([]M, error)
	// FindAllPaginated returns a paginated list of all records.
	FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (*pagination.Page[M], error)
	// FindAllPaginatedBy returns a paginated list of records filtered by a condition.
	FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where Where, preloads ...Preload) (*pagination.Page[M], error)
	// FindAllOrdered returns all records ordered by given sort criteria.
	FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...Preload) ([]M, error)
}

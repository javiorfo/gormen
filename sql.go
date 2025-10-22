package gormen

import (
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/gormen/where"
)

// Preload represents a string identifier for preloading related entities in Gorm
type Preload = string

// Join represents a string identifier for SQL join operations in Gorm
type Join = string

// conditions maps Condition instances to their logical operator type (e.g., AND, OR).
type conditions = map[where.Condition]int

// Where holds SQL query conditions and join clauses to build complex queries.
type Where struct {
	conditions conditions
	joins      []Join
}

// NewWhere creates a new Where instance with an initial condition and no logical operator.
func NewWhere(c where.Condition) *Where {
	return &Where{conditions: conditions{c: types.None}}
}

// Conditions returns the map of conditions and their logical operators.
func (w Where) Conditions() conditions {
	return w.conditions
}

// Joins returns the list of join clauses included in the query.
func (w Where) Joins() []Join {
	return w.joins
}

// And adds a condition combined with a logical AND to the Where clause.
func (w *Where) And(c where.Condition) *Where {
	w.conditions[c] = types.And
	return w
}

// Or adds a condition combined with a logical OR to the Where clause.
func (w *Where) Or(c where.Condition) *Where {
	w.conditions[c] = types.Or
	return w
}

// WithJoin sets the list of join clauses to include in the Where query.
func (w *Where) WithJoin(joins ...Join) *Where {
	w.joins = joins
	return w
}

// Build finalizes and returns a copy of the Where instance.
func (w *Where) Build() Where {
	return *w
}

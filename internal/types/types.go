package types

// Logical operator constants for conditions.
const (
	None = iota // No logical operator
	And         // Logical AND operator
	Or          // Logical OR operator
)

// CRUD operation constants.
const (
	Create = iota // Create operation
	Delete        // Delete operation
	Save          // Save (create or update) operation
)

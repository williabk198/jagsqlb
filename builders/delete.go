package builders

import incondition "github.com/williabk198/jagsqlb/internal/condition"

// DeleteBuilder defines the functions needed to build an SQL "DELETE" statement.
type DeleteBuilder interface {
	ReturningBuilder

	// Using defines a table to be referenced within the "DELETE" query.
	// This is fundamentally equivalent to using "FROM" in an "SELECT" statement.
	Using(table string) DeleteBuilder

	// Where sets the conditions for which items will be deleted from the database.
	// If `moreConditions` is not nil, the provided conditions will be concatenated with "AND" when built.
	Where(condition incondition.Condition, moreConditions ...incondition.Condition) ReturningWhereBuilder
}

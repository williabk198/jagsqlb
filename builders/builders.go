package builders

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
)

type Builder interface {
	// Build collates the data of the underlying implementation and returns the raw query as a string,
	// query parameters as a slice of any type a or an error if one was encountered
	Build() (query string, queryParams []any, err error)
}

type WhereBuilder[T any] interface {
	// And creates additional conditions that are conjoined using "AND"
	And(incondition.Condition, ...incondition.Condition) T
	// Or creates additional conditions that are conjoined using "OR"
	Or(incondition.Condition, ...incondition.Condition) T
}

type ReturningWhereBuilder interface {
	ReturningBuilder
	WhereBuilder[ReturningWhereBuilder]
}

type ReturningBuilder interface {
	Builder

	// Returning sets what columns to return
	Returning(column string, moreColumns ...string) Builder
}

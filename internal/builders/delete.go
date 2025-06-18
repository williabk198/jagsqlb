package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
)

type deleteBuilder struct{}

// Build implements builders.DeleteBuilder.
func (d deleteBuilder) Build() (query string, queryParams []any, err error) {
	panic("unimplemented")
}

// Using implements builders.DeleteBuilder.
func (d deleteBuilder) Using(tableName string) builders.DeleteBuilder {
	panic("unimplemented")
}

// Where implements builders.DeleteBuilder.
func (d deleteBuilder) Where(incondition.Condition, ...incondition.Condition) builders.ReturningWhereBuilder {
	panic("unimplemented")
}

// Returning implements builders.DeleteBuilder.
func (d deleteBuilder) Returning(column string, moreColumns ...string) builders.Builder {
	panic("unimplemented")
}

func NewDeleteBuilder(table string) builders.DeleteBuilder {
	panic("unimplemented")
}

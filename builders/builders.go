package builders

import (
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type SelectBuilder interface {
	Builder
	Table(table string, columns ...string) SelectBuilder
	Where(intypes.Condition, ...intypes.Condition) WhereBuilder
}

type WhereBuilder interface {
	Builder
	And(intypes.Condition, ...intypes.Condition) WhereBuilder
	Or(intypes.Condition, ...intypes.Condition) WhereBuilder
}

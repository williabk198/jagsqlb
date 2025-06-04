package builders

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	"github.com/williabk198/jagsqlb/types"
)

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type SelectBuilder interface {
	OrderByPaginationBuilders
	Table(table string, columns ...string) SelectBuilder
	Where(inconds.Condition, ...inconds.Condition) WhereBuilder
}

type WhereBuilder interface {
	OrderByPaginationBuilders
	And(inconds.Condition, ...inconds.Condition) WhereBuilder
	Or(inconds.Condition, ...inconds.Condition) WhereBuilder
}

type OrderByPaginationBuilders interface {
	OrderByBuilder
	OrderBy(types.ColumnOrdering, ...types.ColumnOrdering) OrderByBuilder
}

type OrderByBuilder interface {
	OffsetBuilder
	Offset(uint) OffsetBuilder
}

type OffsetBuilder interface {
	Builder
	Limit(uint) Builder
}

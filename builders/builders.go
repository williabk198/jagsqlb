package builders

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	"github.com/williabk198/jagsqlb/types"
)

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type SelectBuilder interface {
	OrderByPaginationBuilders
	Table(table string, columns ...string) SelectBuilder
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder
	Where(inconds.Condition, ...inconds.Condition) WhereBuilder
}

type JoinBuilder interface {
	OrderByPaginationBuilders
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder
	Where(condition inconds.Condition, moreConditions ...inconds.Condition) WhereBuilder
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

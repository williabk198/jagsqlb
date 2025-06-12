package builders

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
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
	Where(incondition.Condition, ...incondition.Condition) WhereBuilder
}

type JoinBuilder interface {
	OrderByPaginationBuilders
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder
	Where(condition incondition.Condition, moreConditions ...incondition.Condition) WhereBuilder
}

type WhereBuilder interface {
	OrderByPaginationBuilders
	And(incondition.Condition, ...incondition.Condition) WhereBuilder
	Or(incondition.Condition, ...incondition.Condition) WhereBuilder
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

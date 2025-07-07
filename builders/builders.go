package builders

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	"github.com/williabk198/jagsqlb/types"
)

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type DeleteBuilder interface {
	ReturningBuilder
	Using(table string) DeleteBuilder
	Where(condition incondition.Condition, moreConditions ...incondition.Condition) ReturningWhereBuilder
}

type SelectBuilder interface {
	OrderByPaginationBuilders
	Table(table string, columns ...string) SelectBuilder
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder
	Where(incondition.Condition, ...incondition.Condition) SelectWhereBuilder
}

type JoinBuilder interface {
	OrderByPaginationBuilders
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder
	Where(condition incondition.Condition, moreConditions ...incondition.Condition) SelectWhereBuilder
}

type WhereBuilder[T any] interface {
	And(incondition.Condition, ...incondition.Condition) T
	Or(incondition.Condition, ...incondition.Condition) T
}

type ReturningWhereBuilder interface {
	ReturningBuilder
	WhereBuilder[ReturningWhereBuilder]
}

type SelectWhereBuilder interface {
	OrderByPaginationBuilders
	WhereBuilder[SelectWhereBuilder]
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

type ReturningBuilder interface {
	Builder
	Returning(column string, moreColumns ...string) Builder
}

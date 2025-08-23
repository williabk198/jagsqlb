package builders

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	"github.com/williabk198/jagsqlb/types"
)

type SelectBuilder interface {
	OrderByPaginationBuilders
	// Table adds an addition table to select from as well as any columns that should be returned in the result set.
	Table(table string, columns ...string) SelectBuilder
	// Join defines a join clause to be used in the "SELECT" query as well as any columns from the joining table to add to the result set.
	//
	// For Example:
	//
	//    query, _, err := jagsqlb.NewSqlBuilder().Select("table1 AS t1", "col1").Join(
	//        join.TypeInner,
	//        "table2 AS t2",
	//        join.On(condition.Equal("t1.col1", condition.ColumnValue("t2.col2"))),
	//        "col3", "col4",
	//    ).Build()
	//
	// Will result in:
	//
	//    query = `SELECT "t1"."col1", "t2"."col3", "t2"."col4" FROM "table1" AS "t1" INNER JOIN "table2" AS "t2" ON "t1"."col1" = "t2"."col2";`
	//    err = nil
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder

	// Where sets the conditions for which items will be selected from the database.
	// If `moreConditions` is not nil, the provided conditions will be concatenated with "AND" when built.
	Where(incondition.Condition, ...incondition.Condition) SelectWhereBuilder
}

type JoinBuilder interface {
	OrderByPaginationBuilders

	// Join defines a join clause to be used in the "SELECT" query as well as any columns from the joining table to add to the result set.
	Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) JoinBuilder

	// Where sets the conditions for which items will be selected from the database.
	// If `moreConditions` is not nil, the provided conditions will be concatenated with "AND" when built.
	Where(condition incondition.Condition, moreConditions ...incondition.Condition) SelectWhereBuilder
}

type SelectWhereBuilder interface {
	OrderByPaginationBuilders
	WhereBuilder[SelectWhereBuilder]
}

type OrderByPaginationBuilders interface {
	OffsetBuilder

	// OrderBy sets what columns to sort the result set by
	OrderBy(types.ColumnOrdering, ...types.ColumnOrdering) OffsetBuilder
}

type OffsetBuilder interface {
	LimitBuilder

	// Offset sets how far into the result set
	Offset(uint) LimitBuilder
}

type LimitBuilder interface {
	Builder
	// Limit sets how many items will be in the result set
	Limit(uint) Builder
}

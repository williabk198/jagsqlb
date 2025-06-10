package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	"github.com/williabk198/jagsqlb/types"
)

type joinCondition struct {
	joinTable    intypes.Table
	joinType     injoin.Type
	joinRelation injoin.Relation
}

type joinBuilder struct {
	selectBuilder selectBuilder
	joins         []joinCondition
	errs          intypes.ErrorSlice
}

func (jb joinBuilder) Build() (query string, queryParams []any, err error) {
	// TODO: check `jb.errs` first before attempting the rest of the build
	panic("not implemented") // TODO: Implement
}

func (jb joinBuilder) Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) builders.JoinBuilder {
	// TODO: parse the table string and any provided columns and accumulate any errors into `jb.errss`
	panic("not implemented") // TODO: Implement
}

func (jb joinBuilder) Where(condition inconds.Condition, moreConditions ...inconds.Condition) builders.WhereBuilder {
	panic("not implemented") // TODO: Implement
}

func (jb joinBuilder) Limit(limit uint) builders.Builder {
	panic("not implemented") // TODO: Implement
}

func (jb joinBuilder) Offset(offset uint) builders.OffsetBuilder {
	panic("not implemented") // TODO: Implement
}

func (jb joinBuilder) OrderBy(ordering types.ColumnOrdering, moreOrderings ...types.ColumnOrdering) builders.OrderByBuilder {
	panic("not implemented") // TODO: Implement
}

package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type insertBuilder struct {
	table   intypes.Table
	columns []intypes.Column
	values  [][]any

	errs intypes.ErrorSlice
}

func (ib insertBuilder) Build() (query string, params []any, err error) {
	panic("not implemented")
}

func (ib insertBuilder) Values(vals []any, moreVals ...[]any) builders.ReturningBuilder {
	panic("not implemented") // TODO: Implement
}

func (ib insertBuilder) Columns(column string, moreColumns ...string) builders.InsertValueBuilder {
	panic("not implemented") // TODO: Implement
}

func (ib insertBuilder) DefaultValues() builders.ReturningBuilder {
	panic("not implemented") // TODO: Implement
}

func (ib insertBuilder) Data(data any, moreData ...any) builders.ReturningBuilder {
	panic("not implemented") // TODO: Implement
}

func NewInsertBuilder(table string) builders.InsertBuilder {
	panic("not implemented") // TODO: Implement
}

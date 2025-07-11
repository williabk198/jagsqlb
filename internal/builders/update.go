package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type updateBuilder struct {
	table     intypes.Table
	columns   []intypes.Column
	vals      []any
	fromTable []intypes.Table
	errs      intypes.ErrorSlice
}

// Build implements builders.UpdateBuilder.
func (u updateBuilder) Build() (query string, queryParams []any, err error) {
	panic("unimplemented")
}

// SetMap implements builders.UpdateBuilder.
func (u updateBuilder) SetMap(colValMap map[string]any) builders.UpdateFromBuilder {
	panic("unimplemented")
}

// SetStruct implements builders.UpdateBuilder.
func (u updateBuilder) SetStruct(value any) builders.UpdateFromBuilder {
	panic("unimplemented")
}

// From implements builders.UpdateBuilder.
func (u updateBuilder) From(table string, moreTables ...string) builders.ReturningWhereBuilder {
	panic("unimplemented")
}

func NewUpdateBuilder(table string) builders.UpdateBuilder {
	panic("unimplemented")
}

package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type whereBuilder struct {
	mainQuery builders.Builder
}

func (w whereBuilder) Build() (query string, queryParams []any, err error) {
	panic("unimplemented")
}

func (w whereBuilder) And(intypes.Condition, ...intypes.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

func (w whereBuilder) Or(intypes.Condition, ...intypes.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

func NewWhereBuilder() builders.WhereBuilder {
	return whereBuilder{}
}

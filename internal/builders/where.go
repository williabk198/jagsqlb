package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type whereBuilder struct {
	mainQuery  builders.Builder
	conditions []whereCondition
}

func (w whereBuilder) Build() (query string, queryParams []any, err error) {
	panic("unimplemented")
}

func (w whereBuilder) And(cond intypes.Condition, additionalConds ...intypes.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

func (w whereBuilder) Or(cond intypes.Condition, additionalConds ...intypes.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

type whereCondition struct {
	conjunction string
	condition   intypes.Condition
}

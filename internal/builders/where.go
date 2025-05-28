package inbuilders

import (
	"github.com/williabk198/jagsqlb/builders"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
)

type whereBuilder struct {
	mainQuery  builders.Builder
	conditions []whereCondition
}

func (w whereBuilder) Build() (query string, queryParams []any, err error) {
	panic("unimplemented")
}

func (w whereBuilder) And(cond inconds.Condition, additionalConds ...inconds.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

func (w whereBuilder) Or(cond inconds.Condition, additionalConds ...inconds.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

type whereCondition struct {
	conjunction string
	condition   inconds.Condition
}

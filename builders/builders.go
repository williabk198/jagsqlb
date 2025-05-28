package builders

import inconds "github.com/williabk198/jagsqlb/internal/conditions"

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type SelectBuilder interface {
	Builder
	Table(table string, columns ...string) SelectBuilder
	Where(inconds.Condition, ...inconds.Condition) WhereBuilder
}

type WhereBuilder interface {
	Builder
	And(inconds.Condition, ...inconds.Condition) WhereBuilder
	Or(inconds.Condition, ...inconds.Condition) WhereBuilder
}

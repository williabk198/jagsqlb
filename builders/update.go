package builders

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
)

type UpdateBuilder interface {
	UpdateFromBuilder
	SetMap(colValMap map[string]any) UpdateFromWhereBuilder
	SetStruct(value any) UpdateFromWhereBuilder
}

type UpdateFromBuilder interface {
	Builder
	From(table string, moreTable ...string) ReturningWhereBuilder
}

type UpdateFromWhereBuilder interface {
	UpdateFromBuilder
	Where(cond incondition.Condition, moreConds ...incondition.Condition) ReturningWhereBuilder
}

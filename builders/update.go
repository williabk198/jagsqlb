package builders

type UpdateBuilder interface {
	UpdateFromBuilder
	SetMap(colValMap map[string]any) UpdateFromBuilder
	SetStruct(value any) UpdateFromBuilder
}

type UpdateFromBuilder interface {
	Builder
	From(table string, moreTable ...string) ReturningWhereBuilder
}

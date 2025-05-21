package builders

type Builder interface {
	Build() (query string, queryParams []any, err error)
}

type SelectBuilder interface {
	Builder
	Table(table string, columns ...string) SelectBuilder
}

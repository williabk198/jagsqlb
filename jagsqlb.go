package jagsqlb

import (
	"github.com/williabk198/jagsqlb/builders"
	inbuilders "github.com/williabk198/jagsqlb/internal/builders"
)

// SqlBuilder defines the operations of an SQL builder
type SqlBuilder interface {
	Delete(table string) builders.DeleteBuilder
	Select(table string, columns ...string) builders.SelectBuilder
	Update(table string) builders.UpdateBuilder
}

type sqlBuilder struct{}

func (sb sqlBuilder) Delete(table string) builders.DeleteBuilder {
	return inbuilders.NewDeleteBuilder(table)
}

func (sb sqlBuilder) Insert(table string) builders.InsertBuilder {
	return inbuilders.NewInsertBuilder(table)
}

func (sb sqlBuilder) Select(table string, columns ...string) builders.SelectBuilder {
	return inbuilders.NewSelectBuilder(table, columns...)
}

func (sb sqlBuilder) Update(table string) builders.UpdateBuilder {
	return inbuilders.NewUpdateBuilder(table)
}

// NewSelectBuilder creates and returns a reusable SQL Builder
func NewSqlBuilder() SqlBuilder {
	return sqlBuilder{}
}

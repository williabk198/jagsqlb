package jagsqlb

import (
	"github.com/williabk198/jagsqlb/builders"
	inbuilders "github.com/williabk198/jagsqlb/internal/builders"
)

type SqlBuilder interface {
	Select(table string, columns ...string) builders.SelectBuilder
}

type sqlBuilder struct{}

func (sb sqlBuilder) Select(table string, columns ...string) builders.SelectBuilder {
	return inbuilders.NewSelectBuilder(table, columns...)
}

func NewSqlBuilder() SqlBuilder {
	return sqlBuilder{}
}

package inbuilders

import (
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	inutilities "github.com/williabk198/jagsqlb/internal/utilities"
)

type selectBuilder struct {
	tables  []intypes.Table
	columns []intypes.SelectColumn
	errs    intypes.ErrorSlice
}

func (s selectBuilder) Build() (query string, params []any, err error) {
	if len(s.errs) > 0 {
		return "", nil, s.errs
	}

	var columnStr string
	if len(s.tables) == 1 && len(s.columns) > 0 {
		// If there is only one table defined, we don't need the table prefixes that you'd get by using
		// `inutilities.JoinData`. So, just get the column names
		columnStr = inutilities.CoalesceSelectColumnNamesString(s.columns)

	} else if len(s.columns) > 0 {
		columnStr = inutilities.CoalesceSelectColumnsFullString(s.columns)
	}

	sb := new(strings.Builder)
	sb.WriteString("SELECT ")
	sb.WriteString(columnStr)
	sb.WriteString("FROM ")
	sb.WriteString(inutilities.CoalesceTablesString(s.tables))
	sb.WriteRune(';')

	return sb.String(), nil, nil
}

func (s selectBuilder) Table(table string, columns ...string) builders.SelectBuilder {
	parsedTable, err := tableParser.Parse(table)
	if err != nil {
		s.errs = append(s.errs, err)
	}
	s.tables = append(s.tables, parsedTable)

	for _, col := range columns {
		parsedColumn, err := selectColumnParser.Parse(col)
		if err != nil {
			s.errs = append(s.errs, err)
		}
		parsedColumn.Table = &parsedTable
		s.columns = append(s.columns, parsedColumn)
	}

	return s
}

func (s selectBuilder) Where(cond inconds.Condition, additionalConds ...inconds.Condition) builders.WhereBuilder {
	panic("unimplemented")
}

func NewSelectBuilder(table string, columns ...string) builders.SelectBuilder {
	sbuilder := selectBuilder{}
	return sbuilder.Table(table, columns...)
}

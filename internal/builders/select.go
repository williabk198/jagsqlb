package inbuilders

import (
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	inutilities "github.com/williabk198/jagsqlb/internal/utilities"
	"github.com/williabk198/jagsqlb/types"
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
		// `inutilities.CoalesceSelectColumnsFullString`. So, just get the column names
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

func (s selectBuilder) Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) builders.JoinBuilder {
	jb := joinBuilder{
		selectBuilder: s,
	}
	return jb.Join(joinType, table, joinRelation, includeColumns...)
}

func (s selectBuilder) Where(cond incondition.Condition, additionalConds ...incondition.Condition) builders.SelectWhereBuilder {
	var wb builders.SelectWhereBuilder = selectWhereBuilder{
		mainQuery: s,
		conditions: []whereCondition{
			{condition: cond},
		},
	}

	if len(additionalConds) > 0 {
		wb = wb.And(additionalConds[0], additionalConds[1:]...)
	}

	return wb
}

// Limit implements builders.SelectBuilder.
func (s selectBuilder) Limit(limit uint) builders.Builder {
	return limitBuilder{
		precedingBuilder: s,
		limit:            limit,
	}
}

// Offset implements builders.SelectBuilder.
func (s selectBuilder) Offset(offset uint) builders.LimitBuilder {
	return offsetBuilder{
		precedingBuilder: s,
		offset:           offset,
	}
}

// OrderBy implements builders.SelectBuilder.
func (s selectBuilder) OrderBy(columnOrder types.ColumnOrdering, moreColumnOrders ...types.ColumnOrdering) builders.OffsetBuilder {
	return orderByBuilder{
		precedingBuilder: s,
		columnOrderings:  append([]types.ColumnOrdering{columnOrder}, moreColumnOrders...),
	}
}

func NewSelectBuilder(table string, columns ...string) builders.SelectBuilder {
	sbuilder := selectBuilder{}
	return sbuilder.Table(table, columns...)
}

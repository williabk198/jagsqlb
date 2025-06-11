package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	inutilities "github.com/williabk198/jagsqlb/internal/utilities"
	"github.com/williabk198/jagsqlb/types"
)

type joinCondition struct {
	joinTable    intypes.Table
	joinType     injoin.Type
	joinRelation injoin.Relation
}

type joinBuilder struct {
	selectBuilder selectBuilder
	joins         []joinCondition
	errs          intypes.ErrorSlice
}

func (jb joinBuilder) Build() (query string, queryParams []any, err error) {
	if len(jb.errs) > 0 {
		return "", nil, jb.errs
	}

	sb := new(strings.Builder)

	// Need to build the select query manually here since `selectBuilder.Build` doesn't produce
	// the desired string. Mainly, it won't prepend table data if only one table was defined
	// in `selectBuilder`
	columnStr := inutilities.CoalesceSelectColumnsFullString(jb.selectBuilder.columns)
	tableStr := inutilities.CoalesceTablesString(jb.selectBuilder.tables)

	sb.WriteString("SELECT ")
	sb.WriteString(columnStr)
	sb.WriteString("FROM ")
	sb.WriteString(tableStr)

	for _, joinCond := range jb.joins {
		sb.WriteRune(' ')
		sb.WriteString(string(joinCond.joinType))
		sb.WriteRune(' ')
		sb.WriteString(joinCond.joinTable.String())
		sb.WriteRune(' ')
		sb.WriteString(joinCond.joinRelation.Keyword)
		sb.WriteRune(' ')

		if columnStr, ok := joinCond.joinRelation.Relation.(string); ok && joinCond.joinRelation.Keyword == "USING" {
			column, err := columnParser.Parse(columnStr)
			if err != nil {
				return "", nil, fmt.Errorf("USING column %q was malformed: %w", columnStr, err)
			}
			sb.WriteRune('(')
			sb.WriteString(column.String())
			sb.WriteRune(')')
			continue
		}

		if conditions, ok := joinCond.joinRelation.Relation.([]inconds.Condition); ok && joinCond.joinRelation.Keyword == "ON" {
			condStr, condParams, err := conditions[0].Parameterize()
			if err != nil {
				return "", nil, fmt.Errorf("failed to parameterize ON condition for %q: %w", joinCond.joinType, err)
			}
			queryParams = append(queryParams, condParams...)
			sb.WriteString(condStr)

			for i := 1; i < len(conditions); i++ {
				sb.WriteString(" AND ")
				condStr, condParams, err = conditions[i].Parameterize()
				if err != nil {
					return "", nil, fmt.Errorf("failed to parameterize ON condition for %q: %w", joinCond.joinType, err)
				}

				queryParams = append(queryParams, condParams...)
				sb.WriteString(condStr)
			}
			continue
		}

		return "", nil, fmt.Errorf("invalid join relation type(%T) with %q keyword", joinCond.joinRelation.Relation, joinCond.joinRelation.Keyword)
	}
	sb.WriteRune(';')

	return sb.String(), queryParams, nil
}

func (jb joinBuilder) Join(joinType injoin.Type, table string, joinRelation injoin.Relation, includeColumns ...string) builders.JoinBuilder {
	tableData, err := tableParser.Parse(table)
	if err != nil {
		jb.errs = append(jb.errs, fmt.Errorf("failed to parse table in JOIN clause: %w", err))
		return jb
	}

	for _, colStr := range includeColumns {
		columnData, err := selectColumnParser.Parse(colStr)
		if err != nil {
			jb.errs = append(jb.errs, fmt.Errorf("failed to parse column %q in %s of %s: %w", colStr, joinType, tableData.Name, err))
			return jb
		}
		columnData.Table = &tableData
		jb.selectBuilder.columns = append(jb.selectBuilder.columns, columnData)
	}

	jb.joins = append(jb.joins, joinCondition{
		joinTable:    tableData,
		joinType:     joinType,
		joinRelation: joinRelation,
	})
	return jb
}

func (jb joinBuilder) Where(condition inconds.Condition, moreConditions ...inconds.Condition) builders.WhereBuilder {
	var wb builders.WhereBuilder = whereBuilder{
		mainQuery: jb,
		conditions: []whereCondition{
			{condition: condition},
		},
	}

	if len(moreConditions) > 0 {
		wb = wb.And(moreConditions[0], moreConditions[1:]...)
	}

	return wb
}

func (jb joinBuilder) Limit(limit uint) builders.Builder {
	return limitBuilder{
		precedingBuilder: jb,
		limit:            limit,
	}
}

func (jb joinBuilder) Offset(offset uint) builders.OffsetBuilder {
	return offsetBuilder{
		precedingBuilder: jb,
		offset:           offset,
	}
}

func (jb joinBuilder) OrderBy(ordering types.ColumnOrdering, moreOrderings ...types.ColumnOrdering) builders.OrderByBuilder {
	return orderByBuilder{
		precedingBuilder: jb,
		columnOrderings:  append([]types.ColumnOrdering{ordering}, moreOrderings...),
	}
}

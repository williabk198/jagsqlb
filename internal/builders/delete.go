package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type deleteBuilder struct {
	table       string
	usingTables []intypes.Table
	errs        intypes.ErrorSlice
}

// Build implements builders.DeleteBuilder.
func (d deleteBuilder) Build() (query string, queryParams []any, err error) {
	if len(d.errs) > 0 {
		return "", nil, d.errs
	}

	table, err := tableParser.Parse(d.table)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse table data in base delete builder: %w", err)
	}

	sb := new(strings.Builder)
	sb.WriteString("DELETE FROM ")
	sb.WriteString(table.String())

	if len(d.usingTables) == 0 {
		sb.WriteRune(';')
		return sb.String(), nil, nil
	}

	sb.WriteString(" USING ")
	sb.WriteString(d.usingTables[0].String())

	for i := 1; i < len(d.usingTables); i++ {
		sb.WriteString(", ")
		sb.WriteString(d.usingTables[i].String())
	}

	sb.WriteRune(';')
	return sb.String(), nil, nil
}

// Using implements builders.DeleteBuilder.
func (d deleteBuilder) Using(table string) builders.DeleteBuilder {
	tableData, err := tableParser.Parse(table)
	if err != nil {
		d.errs = append(d.errs, err)
		return d
	}

	d.usingTables = append(d.usingTables, tableData)
	return d
}

// Where implements builders.DeleteBuilder.
func (d deleteBuilder) Where(condition incondition.Condition, moreConditions ...incondition.Condition) builders.ReturningWhereBuilder {
	var rwb builders.ReturningWhereBuilder = returningWhereBuilder{
		mainQuery: d,
		conditions: whereConditions{
			{condition: condition},
		},
	}

	if len(moreConditions) > 0 {
		rwb = rwb.And(moreConditions[0], moreConditions[1:]...)
	}

	return rwb
}

// Returning implements builders.DeleteBuilder.
func (d deleteBuilder) Returning(column string, moreColumns ...string) builders.Builder {
	rb := returningBuilder{
		prevBuilder: d,
	}
	return rb.Returning(column, moreColumns...)
}

func NewDeleteBuilder(table string) builders.DeleteBuilder {
	return deleteBuilder{
		table: table,
	}
}

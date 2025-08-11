package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	"github.com/williabk198/jagsqlb/internal/utilities/parsers"
)

type updateBuilder struct {
	table      intypes.Table
	columns    []intypes.Column
	vals       []any
	fromTables []intypes.Table
	errs       intypes.ErrorSlice
}

// Build implements builders.UpdateBuilder.
func (u updateBuilder) Build() (query string, queryParams []any, err error) {
	if len(u.errs) > 0 {
		return "", nil, fmt.Errorf("failed to build base update query: %w", u.errs)
	}

	sb := new(strings.Builder)
	sb.WriteString("UPDATE ")
	sb.WriteString(u.table.String())

	sb.WriteString(" SET ")
	sb.WriteString(u.columns[0].String())
	sb.WriteRune('=')
	if cv, ok := u.vals[0].(incondition.ColumnValue); ok {
		sb.WriteString(cv.ColumnName)
	} else {
		sb.WriteRune('?')
	}

	for i := 1; i < len(u.columns); i++ {
		sb.WriteString(", ")
		sb.WriteString(u.columns[i].String())
		sb.WriteRune('=')
		if cv, ok := u.vals[i].(incondition.ColumnValue); ok {
			sb.WriteString(cv.ColumnName)
		} else {
			sb.WriteRune('?')
		}
	}

	if len(u.fromTables) > 0 {
		sb.WriteString(" FROM ")
		sb.WriteString(u.fromTables[0].String())
		for i := 1; i < len(u.fromTables); i++ {
			sb.WriteString(", ")
			sb.WriteString(u.fromTables[i].String())
		}
	}

	sb.WriteRune(';')

	// Go through the values and remove any value that is of the type `incondition.ColumnValue`.
	// Doing this here instead of inside the previous maybe slightly more inefficient, but easier to understand.
	// Performance issues will be handled later, if the need arises.
	for i := range u.vals {
		if _, ok := u.vals[i].(incondition.ColumnValue); ok {
			u.vals = append(u.vals[:i], u.vals[i+1:])
		}
	}

	return finalizeQuery(sb.String(), 0), u.vals, nil
}

// SetMap implements builders.UpdateBuilder.
func (u updateBuilder) SetMap(colValMap map[string]any) builders.UpdateFromWhereBuilder {
	u.columns = make([]intypes.Column, len(colValMap))
	u.vals = make([]any, len(colValMap))

	i := 0
	for k, v := range colValMap {
		colData, err := columnParser.Parse(k)
		if err != nil {
			u.errs = append(u.errs, err)
			return u
		}
		u.columns[i] = colData
		u.vals[i] = v
		i++
	}

	return u
}

// SetStruct implements builders.UpdateBuilder.
func (u updateBuilder) SetStruct(value any) builders.UpdateFromWhereBuilder {
	cols, vals, err := parsers.ParseColumnTag(value)
	if err != nil {
		u.errs = append(u.errs, fmt.Errorf("failed to process argument of SetStruct: %w", err))
		return u
	}

	u.columns = make([]intypes.Column, len(cols))
	for i, c := range cols {
		colData, err := columnParser.Parse(c)
		if err != nil {
			u.errs = append(u.errs, err)
			return u
		}
		u.columns[i] = colData
	}
	u.vals = vals

	return u
}

// From implements builders.UpdateBuilder.
func (u updateBuilder) From(table string, moreTables ...string) builders.ReturningWhereBuilder {
	tableData, err := tableParser.Parse(table)
	if err != nil {
		u.errs = append(u.errs, err)
		return returningWhereBuilder{
			mainQuery: u,
		}
	}
	u.fromTables = append(u.fromTables, tableData)

	for _, mt := range moreTables {
		tableData, err = tableParser.Parse(mt)
		if err != nil {
			u.errs = append(u.errs, err)
			return returningWhereBuilder{
				mainQuery: u,
			}
		}
		u.fromTables = append(u.fromTables, tableData)
	}

	return returningWhereBuilder{
		mainQuery: u,
	}
}

func (u updateBuilder) Where(cond incondition.Condition, moreConds ...incondition.Condition) builders.ReturningWhereBuilder {
	rwb := returningWhereBuilder{
		mainQuery:      u,
		existingParams: len(u.vals),
	}
	return rwb.And(cond, moreConds...)
}

func NewUpdateBuilder(table string) builders.UpdateBuilder {
	ub := updateBuilder{}

	tableData, err := tableParser.Parse(table)
	if err != nil {
		ub.errs = append(ub.errs, err)
		return ub
	}
	ub.table = tableData
	return ub
}

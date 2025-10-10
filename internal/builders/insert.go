package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	"github.com/williabk198/jagsqlb/internal/utilities/parsers"
)

type insertBuilder struct {
	table   intypes.Table
	columns []intypes.Column
	values  [][]any

	errs intypes.ErrorSlice
}

func (ib insertBuilder) Build() (query string, params []any, err error) {
	if len(ib.errs) > 0 {
		return "", nil, fmt.Errorf("error(s) exist preceding the build process of the insert statement: %w", ib.errs)
	}

	if len(ib.columns) == 0 && len(ib.values) == 0 {
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES;", ib.table), nil, nil
	}

	sb := new(strings.Builder)
	sb.WriteString("INSERT INTO ")
	sb.WriteString(ib.table.String())

	if len(ib.columns) > 0 {
		sb.WriteString(" (")
		sb.WriteString(ib.columns[0].String())
		for i := 1; i < len(ib.columns); i++ {
			sb.WriteString(", ")
			sb.WriteString(ib.columns[i].String())
		}
		sb.WriteString(")")
	}

	sb.WriteString(" VALUES")

	for i, val := range ib.values {
		sb.WriteString(" (")
		fmt.Fprintf(sb, "$%d", i*len(val)+1)
		for j := 1; j < len(val); j++ {
			sb.WriteString(", ")
			fmt.Fprintf(sb, "$%d", i*len(val)+j+1)
		}
		sb.WriteRune(')')
	}
	sb.WriteRune(';')

	for _, val := range ib.values {
		params = append(params, val...)
	}

	return sb.String(), params, nil
}

func (ib insertBuilder) Values(vals []any, moreVals ...[]any) builders.ReturningBuilder {
	if len(ib.columns) > 0 && len(ib.columns) != len(vals) {
		ib.errs = append(ib.errs, fmt.Errorf("%d column(s) provided but %d value(s) were given(%v)", len(ib.columns), len(vals), vals))
		return returningBuilder{
			prevBuilder: ib,
		}
	}
	ib.values = append(ib.values, vals)

	for _, mv := range moreVals {
		if len(ib.columns) > 0 && len(ib.columns) != len(mv) {
			ib.errs = append(ib.errs, fmt.Errorf("%d column(s) provided but %d value(s) were given(%v)", len(ib.columns), len(mv), mv))
			return returningBuilder{
				prevBuilder: ib,
			}
		}
		ib.values = append(ib.values, mv)
	}

	return returningBuilder{
		prevBuilder: ib,
	}
}

func (ib insertBuilder) Columns(column string, moreColumns ...string) builders.InsertValueBuilder {
	columnData, err := columnParser.Parse(column)
	if err != nil {
		ib.errs = append(ib.errs, err)
		return ib
	}
	ib.columns = append(ib.columns, columnData)

	for _, mc := range moreColumns {
		columnData, err = columnParser.Parse(mc)
		if err != nil {
			ib.errs = append(ib.errs, err)
			return ib
		}
		ib.columns = append(ib.columns, columnData)
	}
	return ib
}

func (ib insertBuilder) DefaultValues() builders.ReturningBuilder {
	// Ensure that both columns and values are empty
	ib.columns = nil
	ib.values = nil

	return returningBuilder{
		prevBuilder: ib,
	}
}

func (ib insertBuilder) Data(data any, moreData ...any) builders.ReturningBuilder {
	cols, vals, err := parsers.ParseColumnTag(intypes.QueryTypeInsert, data)
	if err != nil {
		ib.errs = append(ib.errs, fmt.Errorf("failed to process argument 0 of Data function: %w", err))
		return returningBuilder{
			prevBuilder: ib,
		}
	}

	moreVals := make([][]any, len(moreData))
	for i, md := range moreData {
		_, valData, err := parsers.ParseColumnTag(intypes.QueryTypeInsert, md)
		if err != nil {
			ib.errs = append(ib.errs, fmt.Errorf("failed to process argument %d of Data function: %w", i+1, err))
			return returningBuilder{
				prevBuilder: ib,
			}
		}
		moreVals[i] = valData
	}

	valBuilder := ib.Columns(cols[0], cols[1:]...)
	return valBuilder.Values(vals, moreVals...)
}

func NewInsertBuilder(table string) builders.InsertBuilder {
	ib := insertBuilder{}
	tableData, err := tableParser.Parse(table)
	if err != nil {
		ib.errs = append(ib.errs, err)
	}
	ib.table = tableData

	return ib
}

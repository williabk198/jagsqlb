package inconds

import (
	"fmt"
	"strings"
)

type SimpleCondition struct {
	ColumnName string
	Operator   string
	Values     []any
}

func (sc SimpleCondition) Parameterize() (string, []any, error) {
	column, err := columnParser.Parse(sc.ColumnName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse column data: %w", err)
	}

	// TODO?: Move to its own function??
	if strings.HasSuffix(sc.Operator, "BETWEEN") {
		sb := new(strings.Builder)
		sb.WriteString(column.String())
		sb.WriteRune(' ')
		sb.WriteString(sc.Operator)
		sb.WriteRune(' ')

		// Check to see if the first value is a ColumnValue
		if columnVal, ok := sc.Values[0].(ColumnValue); ok {
			// If so, then parse out the column data and append it to the builder
			col, err := columnParser.Parse(columnVal.ColumnName)
			if err != nil {
				return "", nil, fmt.Errorf("failed to parse ColumnValue")
			}
			sb.WriteString(col.String())
			sc.Values = sc.Values[1:]
		} else {
			// Otherwise, the value will be parameterized
			sb.WriteString("?")
		}

		sb.WriteString(" AND ")

		// Check to see if the last value is a ColumnValue
		if columnVal, ok := sc.Values[len(sc.Values)-1].(ColumnValue); ok {
			// If so, do the same thing as above
			col, err := columnParser.Parse(columnVal.ColumnName)
			if err != nil {
				return "", nil, fmt.Errorf("failed to parse ColumnValue")
			}
			sb.WriteString(col.String())
			sc.Values = sc.Values[:len(sc.Values)-1]
		} else {
			// Likewise here. Just parameterize the value if it isn't a ColumnValue
			sb.WriteString("?")
		}

		return sb.String(), sc.Values, nil
	}

	// Check to see if the first value is a ColumnValue and is not an "IN" condition
	inOperation := strings.HasSuffix(sc.Operator, "IN")
	if columnVal, ok := sc.Values[0].(ColumnValue); ok && !inOperation {
		// If so, parse out the column data and use it in the returned string
		columnValStr, err := columnParser.Parse(columnVal.ColumnName)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse ColumnValue data: %w", err)
		}
		return fmt.Sprintf("%s %s %s", column, sc.Operator, columnValStr), sc.Values[1:], nil
	}

	// If the slice of values contains a ColumnValue and this is an "IN" condition, then throw an error.
	// The column will be treated as a string; leading to unwanted results.
	if containsColumnValue(sc.Values) && inOperation {
		return "", nil, fmt.Errorf("cannot have a ColumnValue within a parametreized IN condition")
	}
	return fmt.Sprintf("%s %s ?", column, sc.Operator), sc.Values, nil
}

// containsColumnValue takes in a slice of values and checks to see if any are of the type ColumnValue
func containsColumnValue(vals []any) bool {
	for _, val := range vals {
		if _, ok := val.(ColumnValue); ok {
			return true
		}
	}
	return false
}

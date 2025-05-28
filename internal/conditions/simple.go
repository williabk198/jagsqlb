package inconds

import "fmt"

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

	var resultStr string
	switch sc.Operator {
	case "BETWEEN", "NOT BETWEEN":
		resultStr = fmt.Sprintf("%s %s ? AND ?", column, sc.Operator)
	default:
		resultStr = fmt.Sprintf("%s %s ?", column, sc.Operator)
	}

	return resultStr, sc.Values, nil
}

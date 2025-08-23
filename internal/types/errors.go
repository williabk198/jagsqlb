package intypes

import (
	"fmt"
	"strings"
)

var (
	ErrMissingAliasName  = NewInvalidSyntaxError("alias name not provided")
	ErrMissingColumnName = NewInvalidSyntaxError("column name not provided")
	ErrMissingSchemaName = NewInvalidSyntaxError("schema name not provided")
	ErrMissingTableName  = NewInvalidSyntaxError("table name not provided")
)

func NewInvalidSyntaxError(errString string) error {
	return fmt.Errorf("invalid syntax error: %s", errString)
}

type ErrorSlice []error

func (es ErrorSlice) Append(err error) {
	es = append(es, err)
}

func (es ErrorSlice) Error() string {
	if len(es) == 0 {
		return ""
	}

	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("encountered %d error(s)", len(es)))

	for _, err := range es {
		sb.WriteString("\n\t")
		sb.WriteString(err.Error())
	}

	return sb.String()
}

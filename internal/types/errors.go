package intypes

import (
	"fmt"
	"strings"
)

var (
	ErrMissingAliasName  = NewInvalidSytaxError("alias name not provided")
	ErrMissingColumnName = NewInvalidSytaxError("column name not provided")
	ErrMissingSchemaName = NewInvalidSytaxError("schema name not provided")
	ErrMissingTableName  = NewInvalidSytaxError("table name not provided")
)

func NewInvalidSytaxError(errString string) error {
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

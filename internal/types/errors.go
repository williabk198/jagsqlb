package intypes

import "fmt"

var (
	ErrMissingAliasName  = NewInvalidSytaxError("alias name not provided")
	ErrMissingColumnName = NewInvalidSytaxError("column name not provided")
	ErrMissingSchemaName = NewInvalidSytaxError("schema name not provided")
	ErrMissingTableName  = NewInvalidSytaxError("table name not provided")
)

func NewInvalidSytaxError(errString string) error {
	return fmt.Errorf("invalid syntax error: %s", errString)
}

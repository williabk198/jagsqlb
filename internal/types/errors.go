package intypes

import "fmt"

func NewInvalidSytaxError(errString string) error {
	return fmt.Errorf("invalid syntax error: %s", errString)
}

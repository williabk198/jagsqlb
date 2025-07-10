package parsers

import (
	"fmt"
	"reflect"
)

var (
	ErrInputTypeNotStruct = fmt.Errorf("recieved value is not a struct type")
)

// ParseColumnTag expects a struct for `input`. If it isn't a struct, then an error is returned.
// Otherwise, it will look for the `jagsqlb` struct tag which denotes the names of the column in
// the database and returns the mapping of that column to its corresponding value.
func ParseColumnTag(input any) (cols []string, vals []any, err error) {
	// Since we are working with the reflect package, we need to worry about handling panics so that it errors out gracefully,
	// instead of just crashing out.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	inputType := reflect.TypeOf(input)
	inputValue := reflect.ValueOf(input)
	if inputType.Kind() != reflect.Struct {
		return nil, nil, ErrInputTypeNotStruct
	}

	// Since NumField always checks to see if `inputType` is a struct, it's better to just
	// store NumField in a variable so it doesn't have to run that check over and over.
	fieldCount := inputType.NumField()
	cols = make([]string, fieldCount)
	vals = make([]any, fieldCount)

	for i := range fieldCount {
		fieldType := inputType.Field(i)
		fieldVal := inputValue.Field(i)

		if colName := fieldType.Tag.Get("jagsqlb"); colName != "" {
			cols[i] = colName
		} else {
			cols[i] = fieldType.Name
		}
		vals[i] = fieldVal.Interface()
	}

	return cols, vals, nil
}

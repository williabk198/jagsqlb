package parsers

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrInputTypeNotStruct = fmt.Errorf("received value is not a struct type")
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

	cols = []string{}
	vals = []any{}
	for i := range inputType.NumField() {
		fieldType := inputType.Field(i)
		fieldVal := inputValue.Field(i)

		tagVal := fieldType.Tag.Get("jagsqlb")
		splitVals := strings.Split(tagVal, ";")

		tagData := tagData{
			columnName: splitVals[0],
		}
		for i := 1; i < len(splitVals); i++ {
			switch splitVals[i] {
			case "inline":
				tagData.inline = true
			case "omit":
				//TODO?(BW): Consider conditional omits. For example, omit if an insert/update statement, or if empty, and so on.
				tagData.omit = true
			}
		}

		if tagData.omit {
			continue
		}

		if tagData.inline && fieldType.Type.Kind() == reflect.Struct {
			c, v, _ := ParseColumnTag(fieldVal.Interface())
			cols = append(cols, c...)
			vals = append(vals, v...)
			continue
		}

		if tagData.columnName != "" {
			cols = append(cols, tagData.columnName)
		} else {
			cols = append(cols, fieldType.Name)
		}
		vals = append(vals, fieldVal.Interface())
	}

	return cols, vals, nil
}

type tagData struct {
	columnName string
	inline     bool
	omit       bool
}

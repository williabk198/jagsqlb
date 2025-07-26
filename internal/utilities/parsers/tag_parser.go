package parsers

import (
	"fmt"
	"reflect"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
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

		fieldData := fieldVal.Interface()
		if fieldType.Type.Implements(reflect.TypeFor[intypes.QueryMarshaler]()) {
			// If the current field implements the QueryMarshaler interface, the use that
			// to build the value for the column
			qm := fieldVal.Interface().(intypes.QueryMarshaler)
			fieldData, err = qm.MarshalQuery()
			if err != nil {
				return nil, nil, fmt.Errorf(
					"failed to marshal struct data for field %q: %w",
					fieldType.Name, err,
				)
			}
		} else if tagData.inline && fieldType.Type.Kind() == reflect.Struct {
			// Otherwise, if the property is a struct and has been marked as "inline",
			// then recursively call ParseColumnTag to get the columns and values of the nested struct
			c, v, e := ParseColumnTag(fieldVal.Interface())
			if e != nil {
				return nil, nil, fmt.Errorf(
					"failed to marshal nested struct data for field %q: %w",
					fieldType.Name, e,
				)
			}
			cols = append(cols, c...)
			vals = append(vals, v...)
			continue
		}

		if tagData.columnName != "" {
			cols = append(cols, tagData.columnName)
		} else {
			cols = append(cols, fieldType.Name)
		}
		vals = append(vals, fieldData)
	}

	return cols, vals, nil
}

type tagData struct {
	columnName string
	inline     bool
	omit       bool
}

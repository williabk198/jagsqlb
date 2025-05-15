package inutilities

import (
	"fmt"
	"regexp"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

var (
	extraWhitespaceRegex = regexp.MustCompile(`\s{2,}`)
)

func ParseTableData(tableStr string) (intypes.Table, error) {
	// strip out all quotes
	input := strings.ReplaceAll(tableStr, "\"", "")
	input = strings.TrimSpace(input)
	// replace all consecutive whitespace characters with a single space character
	input = extraWhitespaceRegex.ReplaceAllString(input, " ")

	var schema string
	splitInput := strings.Split(input, ".")
	if splitLen := len(splitInput); splitLen == 2 {
		schema = strings.TrimSpace(splitInput[0])
		if schema == "" {
			return intypes.Table{}, intypes.NewInvalidSytaxError(fmt.Sprintf("expected the schema name before the '.' in %q", tableStr))
		}

		input = strings.TrimSpace(splitInput[1])
		if input == "" {
			return intypes.Table{}, intypes.NewInvalidSytaxError(fmt.Sprintf("expected the table name after the '.' in %q", tableStr))
		}

		input = splitInput[1]
	} else if splitLen > 2 {
		return intypes.Table{}, intypes.NewInvalidSytaxError(fmt.Sprintf("too many '.' characters in %q", tableStr))
	}

	// attempt to get table alias by searching for "AS"
	alias, err := getAlias(&input, "AS")

	// if the alias was not found, and there wasn't an error then try to get the alias by looking for the space.
	if err == nil && alias == "" {
		alias, err = getAlias(&input, " ")
	}
	if err != nil {
		return intypes.Table{}, fmt.Errorf("failed to parse table alias in %q: %w", tableStr, err)
	}

	return intypes.Table{
		Alias:  alias,
		Name:   input,
		Schema: schema,
	}, nil
}

func ParseColumnData(input string) (intypes.Column, error) {
	panic("unimplemented")
}

func ParseSelectorColumnData(input string) (intypes.SelectorColumn, error) {
	// TODO: Use ParseColumnData and then parse out the column alias if it exsists
	panic("unimplemented")
}

func getAlias(input *string, seperator string) (string, error) {
	var alias string
	splitInput := strings.Split(*input, seperator)
	if splitLen := len(splitInput); splitLen == 2 {
		*input = strings.TrimSpace(splitInput[0])
		if *input == "" {
			return "", intypes.NewInvalidSytaxError("value before alias definition is empty or whitespace")
		}

		alias = strings.TrimSpace(splitInput[1])
		if alias == "" {
			return "", intypes.NewInvalidSytaxError("alias name not provided")
		}
	} else if splitLen > 2 {
		return "", intypes.NewInvalidSytaxError(fmt.Sprintf("multiple occurences of %q", seperator))
	}

	return alias, nil
}

package inutilities

// TODO: Move this file and associated test to a "parsers" sub-package

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

	if input == "" {
		return intypes.Table{}, intypes.NewInvalidSytaxError("table name not provided")
	}

	return intypes.Table{
		Alias:  alias,
		Name:   input,
		Schema: schema,
	}, nil
}

func ParseColumnData(columnStr string) (intypes.Column, error) {
	// strip out all quotes
	input := strings.ReplaceAll(columnStr, "\"", "")
	input = strings.TrimSpace(input)
	// replace all consecutive whitespace characters with a single space character
	input = extraWhitespaceRegex.ReplaceAllString(input, " ")

	var table *intypes.Table
	lastPeriodIndex := strings.LastIndex(input, ".")
	if lastPeriodIndex != -1 {
		tableStr := input[:lastPeriodIndex]
		t, err := ParseTableData(tableStr)
		if err != nil {
			return intypes.Column{}, fmt.Errorf("failed to parse table data provided in %q: %w", columnStr, err)
		}
		table = &t
		input = strings.TrimSpace(input[lastPeriodIndex+1:])
	}

	// strip out the alias, if one was given
	alias, err := getAlias(&input, "AS")
	if err == nil && alias == "" {
		// if an alias wasn't found using "AS" then try to find the alias using " "
		// the actually result here doesn't matter, just mutating `input` to remove the alias
		getAlias(&input, " ")
	}

	if input == "" {
		return intypes.Column{}, intypes.NewInvalidSytaxError("column name was not provided")
	}

	return intypes.Column{
		Name:  input,
		Table: table,
	}, nil
}

func ParseSelectorColumnData(selectorColumnStr string) (intypes.SelectorColumn, error) {
	// Decided not to use `PraseColumnData` here. That's because it runs `getAlias` to strip
	// out the alias to isolate the column name, and the alias is discarded. Which means a second
	// call to `getAlias` would be needed here to parse out the alias value, and that seemed inefficient.

	// TODO: Refactor. Move the string modifications to its own function since both `ParseColumnData` and
	//       `ParseTableData` use this exact functionality as well

	// strip out all quotes
	input := strings.ReplaceAll(selectorColumnStr, "\"", "")
	input = strings.TrimSpace(input)
	// replace all consecutive whitespace characters with a single space character
	input = extraWhitespaceRegex.ReplaceAllString(input, " ")

	//TODO: Refactor. Move the bellow block to its own function since `ParseColumnData` uses this
	// functionality as well
	var table *intypes.Table
	lastPeriodIndex := strings.LastIndex(input, ".")
	if lastPeriodIndex != -1 {
		tableStr := input[:lastPeriodIndex]
		t, err := ParseTableData(tableStr)
		if err != nil {
			return intypes.SelectorColumn{}, fmt.Errorf("failed to parse table data provided in %q: %w", selectorColumnStr, err)
		}
		table = &t
		input = strings.TrimSpace(input[lastPeriodIndex+1:])
	}

	// attempt to get table alias by searching for "AS"
	alias, err := getAlias(&input, "AS")
	// if the alias was not found, and there wasn't an error then try to get the alias by looking for the space.
	if err == nil && alias == "" {
		alias, err = getAlias(&input, " ")
	}
	if err != nil {
		return intypes.SelectorColumn{}, fmt.Errorf("failed to parse column alias in %q: %w", selectorColumnStr, err)
	}

	return intypes.SelectorColumn{
		Column: intypes.Column{
			Name:  input,
			Table: table,
		},
		Alias: alias,
	}, nil
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

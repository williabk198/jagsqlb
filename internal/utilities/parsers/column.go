package parsers

import (
	"fmt"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type columnParser struct {
}

func (cp columnParser) Parse(columnStr string) (intypes.Column, error) {
	var err error
	var table *intypes.Table
	input := sanitizeInput(columnStr)

	// Get the last index of the "." character.
	lastPeriodIndex := strings.LastIndex(input, ".")
	if lastPeriodIndex != -1 {
		// If a "." character was found, then split the input string.
		tableStr := input[:lastPeriodIndex]                  // the first portion should be related to table information
		input = strings.TrimSpace(input[lastPeriodIndex+1:]) // the latter portion should be related to column data

		table, err = getTableData(&tableStr)
		if err != nil {
			return intypes.Column{}, fmt.Errorf("failed to parse table data provided in %q: %w", columnStr, err)
		}
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

type selectColumnParser struct {
}

func (scp selectColumnParser) Parse(selectColumnStr string) (intypes.SelectColumn, error) {
	// Decided not to use `PraseColumnData` here. That's because it runs `getAlias` to strip
	// out the alias to isolate the column name, and the alias is discarded. Which means a second
	// call to `getAlias` would be needed here to parse out the alias value, and that seemed inefficient.

	var err error
	var table *intypes.Table
	input := sanitizeInput(selectColumnStr)

	lastPeriodIndex := strings.LastIndex(input, ".")
	if lastPeriodIndex != -1 {
		tableStr := input[:lastPeriodIndex]
		table, err = getTableData(&tableStr)
		if err != nil {
			return intypes.SelectColumn{}, fmt.Errorf("failed to parse table data provided in %q: %w", selectColumnStr, err)
		}
		input = strings.TrimSpace(input[lastPeriodIndex+1:])
	}

	// attempt to get table alias by searching for "AS"
	alias, err := getAlias(&input, "AS")
	// if the alias was not found, and there wasn't an error then try to get the alias by looking for the space.
	if err == nil && alias == "" {
		alias, err = getAlias(&input, " ")
	}
	if err != nil {
		return intypes.SelectColumn{}, fmt.Errorf("failed to parse column alias in %q: %w", selectColumnStr, err)
	}

	return intypes.SelectColumn{
		Column: intypes.Column{
			Name:  input,
			Table: table,
		},
		Alias: alias,
	}, nil
}

package parsers

import (
	"fmt"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type tableParser struct{}

func (tp tableParser) Parse(tableStr string) (intypes.Table, error) {

	input := sanitizeInput(tableStr)
	table, err := getTableData(&input)
	if err != nil {
		return intypes.Table{}, fmt.Errorf("failed to parse table data in %q: %w", tableStr, err)
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
	table.Alias = alias
	table.Name = input

	return *table, nil
}

func NewTableParser() Parser[intypes.Table] {
	return tableParser{}
}

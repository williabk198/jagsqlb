package parsers

import (
	"fmt"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type tableParser struct{}

func (tp tableParser) Parse(tableStr string) (intypes.Table, error) {

	remainder := sanitizeInput(tableStr)
	table, err := getTableData(&remainder)
	if err != nil {
		return intypes.Table{}, fmt.Errorf("failed to parse table data from %q: %w", tableStr, err)
	}

	alias, remainder, err := getAlias(remainder)
	if err != nil {
		return intypes.Table{}, fmt.Errorf("failed to parse table alias in %q: %w", tableStr, err)
	}
	table.Alias = alias
	table.Name = remainder

	return *table, nil
}

func NewTableParser() Parser[intypes.Table] {
	return tableParser{}
}

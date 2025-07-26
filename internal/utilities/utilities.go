package inutilities

import (
	"fmt"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

// CoalesceSelectColumnsFullString takes in a slice of SelectColumns and returns them as a comma separated string
// using its fully-qualified definition
func CoalesceSelectColumnsFullString(cols []intypes.SelectColumn) string {
	if len(cols) == 0 {
		return ""
	}

	strSlice := make([]string, len(cols))
	for i, col := range cols {
		strSlice[i] = col.String()
	}

	result := strings.Join(strSlice, ", ")
	if result != "" {
		result += " "
	}

	return result
}

// CoalesceSelectColumnNamesString takes in a slice of SelectColumns and returns them as a comma separated string
// using only the name of the column, and its alias (if one was provided).
func CoalesceSelectColumnNamesString(cols []intypes.SelectColumn) string {
	if len(cols) == 0 {
		return ""
	}

	strSlice := make([]string, len(cols))
	for i, col := range cols {
		if col.Name == "*" {
			strSlice[i] = col.Name
		} else if col.Alias == "" {
			strSlice[i] = fmt.Sprintf("%q", col.Name)
		} else {
			strSlice[i] = fmt.Sprintf("%q AS %q", col.Name, col.Alias)
		}
	}

	result := strings.Join(strSlice, ", ")
	if result != "" {
		result += " "
	}

	return result
}

// CoalesceTablesString takes in a slice of Tables and returns them as a comma separated string using its fully qualified definition
func CoalesceTablesString(tables []intypes.Table) string {
	if len(tables) == 0 {
		return ""
	}

	result := make([]string, len(tables))
	for i, table := range tables {
		result[i] = table.String()
	}

	return strings.Join(result, ", ")
}

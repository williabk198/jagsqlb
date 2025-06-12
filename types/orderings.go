package types

import (
	"fmt"

	"github.com/williabk198/jagsqlb/internal/utilities/parsers"
)

type ordering string

const (
	OrderingAscending  ordering = "ASC"
	OrderingDescending ordering = "DESC"
)

var (
	columnParser = parsers.NewColumnParser()
)

type ColumnOrdering struct {
	ColumnName string
	Ordering   ordering
}

func (co ColumnOrdering) Stringify() (string, error) {
	column, err := columnParser.Parse(co.ColumnName)
	if err != nil {
		return "", fmt.Errorf("failed to parse column name %q: %w", co.ColumnName, err)
	}

	return fmt.Sprintf("%s %s", column.String(), co.Ordering), nil
}

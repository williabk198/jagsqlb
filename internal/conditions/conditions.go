package inconds

import "github.com/williabk198/jagsqlb/internal/utilities/parsers"

var (
	columnParser = parsers.NewColumnParser()
)

// Condition is essentially the same as a Builder, but needed a way to differentiate a Condtion from a Builder
type Condition interface {
	Parameterize() (string, []any, error)
}

// ColumnValue represents a column that will be uses as a value within a condition.
// Meaning, that it facilitates having the following types of conidions: "t1.col1 < t2.col2"
type ColumnValue struct {
	ColumnName string
}

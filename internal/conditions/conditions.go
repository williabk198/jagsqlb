package inconds

import "github.com/williabk198/jagsqlb/internal/utilities/parsers"

var (
	columnParser = parsers.NewColumnParser()
)

// Condition is essentially the same as a Builder, but needed a way to differentiate a Condtion from a Builder
type Condition interface {
	Parameterize() (string, []any, error)
}

// ColumnValue represents a column that will be uses as a value in a condition.
// For example, if you wanted to have a condition like "t1.col1 < t2.col2", then
// you use ColumnValue like so:
//
//	exampleCond := conds.LessThan("t1.col1", conds.ColumnValue{ColumnName: "t2.col2"})
type ColumnValue struct {
	ColumnName string
}

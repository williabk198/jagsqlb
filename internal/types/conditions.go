package intypes

// Condition is essentially the same as a Builder, but needed a way to differentiate a Condtion from a Builder
type Condition interface {
	Parameterize() (string, []any, error)
}

type SimpleCondition struct {
	ColumnName string
	Operator   string
	Values     []any
}

func (sc SimpleCondition) Parameterize() (string, []any, error) {
	// TODO: The returned string should be of the format "<column> <operator> ?"
	panic("unimplemented")
}

type GroupedConditions struct {
	Conjunction string // This value should always be either " AND " or " OR "
	Conditions  []Condition
}

func (gc GroupedConditions) Parameterize() (string, []any, error) {
	// TODO: Loop through and parameterieze each condition. Join the resulting strings using `Conjunction` and with surrounding parenthesis.
	//       Merge returned paramters into single slice, and accumulate any errors into an `ErrorSlice`
	panic("unimplemented")
}

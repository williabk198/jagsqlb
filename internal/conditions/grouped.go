package inconds

type GroupedConditions struct {
	Conjunction string // This value should always be either " AND " or " OR "
	Conditions  []Condition
}

func (gc GroupedConditions) Parameterize() (string, []any, error) {
	// TODO: Loop through and parameterieze each condition. Join the resulting strings using `Conjunction` and with surrounding parenthesis.
	//       Merge returned paramters into single slice, and accumulate any errors into an `ErrorSlice`
	panic("unimplemented")
}

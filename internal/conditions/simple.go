package inconds

type SimpleCondition struct {
	ColumnName string
	Operator   string
	Values     []any
}

func (sc SimpleCondition) Parameterize() (string, []any, error) {
	// TODO: The returned string should be of the format "<column> <operator> ?"
	panic("unimplemented")
}

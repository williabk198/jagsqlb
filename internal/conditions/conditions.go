package inconds

// Condition is essentially the same as a Builder, but needed a way to differentiate a Condtion from a Builder
type Condition interface {
	Parameterize() (string, []any, error)
}

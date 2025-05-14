package intypes

type Table struct {
	Alias  string
	Name   string
	Schema string
}

func (t Table) String() string {
	panic("unimplemented")
}

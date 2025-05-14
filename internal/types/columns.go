package intypes

type Column struct {
	Name  string
	Table Table
}

func (c Column) String() string {
	panic("unimplmented")
}

type SelectorColumn struct {
	Alias string
	Column
}

func (sc SelectorColumn) String() string {
	panic("unimplmented")
}

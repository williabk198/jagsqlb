package intypes

type Table struct {
	Alias  string
	Name   string
	Schema string
}

// ReferenceString returns a string that can be used as a reference to the table.
// If the table has an alias, then this returns said alias. If the table has a schema defined,
// then the both the schema and the table name will be returned in the format "schema.table".
// Otherwise, just the Name field is returned.
func (t Table) ReferenceString() string {
	panic("unimplemented")
}

func (t Table) String() string {
	panic("unimplemented")
}

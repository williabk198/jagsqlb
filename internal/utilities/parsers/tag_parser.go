package parsers

// ColumnTagParser expects a struct for `input`. If it isn't a struct, then an error is returned.
// Otherwise, it will look for the `jagsqlb` struct tag which denotes the names of the column in
// the database and returns the mapping of that column to its corresponding value.
func ColumnTagParser(input any) (map[string]any, error) {
	panic("unimplemented")
}

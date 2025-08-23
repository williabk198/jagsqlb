package builders

type InsertBuilder interface {
	Builder
	InsertValueBuilder

	// Columns defines the list of columns that will be receiving data in the "INSERT" statement
	Columns(column string, moreColumns ...string) InsertValueBuilder

	// Data takes in a struct and an optional slice of structs of the same type and parses out the column and values data.
	//
	// For example:
	//    testData := struct {
	//        PrimaryKey uuid.UUID `jagsqlb:"id;omit"`
	//        Field1     string    `jagsqlb:"field1"`
	//        Field2     int
	//    }{"hello", 42}
	//
	//    query, params, err := jagsqlb.NewSqlBuilder().Insert("table1").Data(testData).Build()
	//
	// Results in the following:
	//
	//    query = `INSERT INTO "table1" ("field1", "Field2") VALUES ($1, $2);`
	//    params = []any{"hello", 42}
	//    err = nil
	Data(data any, moreData ...any) ReturningBuilder

	// DefaultValues will instruct to the database to use the default values for each of the columns in the table
	// instead of providing the values manually.
	DefaultValues() ReturningBuilder
}

type InsertValueBuilder interface {
	// Values associates a value to a column.
	//
	// NOTE: the length of `vals` as well as subsequent entries in `moreVals` must equal the number of columns provided.
	// Meaning, if only two columns were provided, then `vals` and each item in `moreVals` MUST contain exactly two elements.
	Values(vals []any, moreVals ...[]any) ReturningBuilder
}

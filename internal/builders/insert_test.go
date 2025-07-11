package inbuilders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_insertBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		ib        insertBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Single Entry w/ Columns",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}},
				values: [][]any{
					{"test", 13},
				},
			},
			wants: wants{
				query:  `INSERT INTO "table1" ("col1", "col2") VALUES ($1, $2);`,
				params: []any{"test", 13},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Single Entry w/o Columns",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{},
				values: [][]any{
					{"test", 13, true, 1.23},
				},
			},
			wants: wants{
				query:  `INSERT INTO "table1" VALUES ($1, $2, $3, $4);`,
				params: []any{"test", 13, true, 1.23},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Multiple Entries w/ Columns",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}},
				values: [][]any{
					{"test", 13},
					{"more_testing", 97},
				},
			},
			wants: wants{
				query:  `INSERT INTO "table1" ("col1", "col2") VALUES ($1, $2) ($3, $4);`,
				params: []any{"test", 13, "more_testing", 97},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Multiple Entries w/o Columns",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{},
				values: [][]any{
					{"test", 13, true, 1.23},
					{"more_testing", 97, false, 4.56},
				},
			},
			wants: wants{
				query:  `INSERT INTO "table1" VALUES ($1, $2, $3, $4) ($5, $6, $7, $8);`,
				params: []any{"test", 13, true, 1.23, "more_testing", 97, false, 4.56},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Default Values",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			wants: wants{
				query: `INSERT INTO "table1" DEFAULT VALUES;`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; ErrorSlice not Empty",
			ib: insertBuilder{
				errs: intypes.ErrorSlice{assert.AnError},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.ib.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotParams)
		})
	}
}

func Test_insertBuilder_Values(t *testing.T) {
	type args struct {
		vals     []any
		moreVals [][]any
	}

	testInsertBuilder := insertBuilder{
		table:   intypes.Table{Name: "table1"},
		columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}, {Name: "ts"}},
	}

	tests := []struct {
		name string
		ib   insertBuilder
		args args
		want builders.ReturningBuilder
	}{
		{
			name: "Success; Single Value Slice",
			ib:   testInsertBuilder,
			args: args{
				vals: []any{"something", 17, 1.23},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}, {Name: "ts"}},
					values: [][]any{
						{"something", 17, 1.23},
					},
				},
			},
		},
		{
			name: "Success; Multiple Value Slices",
			ib:   testInsertBuilder,
			args: args{
				vals: []any{"something", 17, 1.23},
				moreVals: [][]any{
					{"something_else", 7, 4.56},
				},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}, {Name: "ts"}},
					values: [][]any{
						{"something", 17, 1.23},
						{"something_else", 7, 4.56},
					},
				},
			},
		},
		{
			name: "Error; Incorrect Values Length",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}},
				values:  [][]any{},
			},
			args: args{
				vals: []any{"testing", "too_many_vals"},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "col1"}},
					values:  [][]any{},
					errs: intypes.ErrorSlice{
						fmt.Errorf("1 column(s) provided but 2 value(s) were given(%v)", []any{"testing", "too_many_vals"}),
					},
				},
			},
		},
		{
			name: "Error; Incorrect MoreValues Length",
			ib: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}},
				values:  [][]any{},
			},
			args: args{
				vals: []any{"testing", 56},
				moreVals: [][]any{
					{"testing"},
				},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}},
					values: [][]any{
						{"testing", 56},
					},
					errs: intypes.ErrorSlice{
						fmt.Errorf("2 column(s) provided but 1 value(s) were given(%v)", []any{"testing"}),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ib.Values(tt.args.vals, tt.args.moreVals...))
		})
	}
}

func Test_insertBuilder_Columns(t *testing.T) {
	type args struct {
		column      string
		moreColumns []string
	}
	tests := []struct {
		name string
		ib   insertBuilder
		args args
		want builders.InsertValueBuilder
	}{
		{
			name: "Success; Single Column",
			ib:   insertBuilder{table: intypes.Table{Name: "table1"}},
			args: args{
				column: "col1",
			},
			want: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}},
			},
		},
		{
			name: "Success; Multiple Columns",
			ib:   insertBuilder{table: intypes.Table{Name: "table1"}},
			args: args{
				column:      "col1",
				moreColumns: []string{"col2", "col3"},
			},
			want: insertBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}, {Name: "col2"}, {Name: "col3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ib.Columns(tt.args.column, tt.args.moreColumns...))
		})
	}
}

func Test_insertBuilder_DefaultValues(t *testing.T) {
	tests := []struct {
		name string
		ib   insertBuilder
		want builders.ReturningBuilder
	}{
		{
			name: "Success",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table: intypes.Table{Name: "table1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ib.DefaultValues())
		})
	}
}

func Test_insertBuilder_Data(t *testing.T) {
	type testData struct {
		StrData string `jagsqlb:"string_data"`
		IntData int
	}

	type args struct {
		data     any
		moreData []any
	}

	tests := []struct {
		name string
		ib   insertBuilder
		args args
		want builders.ReturningBuilder
	}{
		{
			name: "Success; Struct with no Tags",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			args: args{
				data: struct{ Data string }{"testing"},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "Data"}},
					values: [][]any{
						{"testing"},
					},
				},
			},
		},
		{
			name: "Success; Struct with Tags",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			args: args{
				data: struct {
					Data int `jagsqlb:"data"`
				}{56},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "data"}},
					values: [][]any{
						{56},
					},
				},
			},
		},
		{
			name: "Success; Multiple Params",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			args: args{
				data: testData{"test1", 42},
				moreData: []any{
					testData{"test2", 93},
				},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "string_data"}, {Name: "IntData"}},
					values: [][]any{
						{"test1", 42},
						{"test2", 93},
					},
				},
			},
		},
		{
			name: "Error; Bad Argument Type",
			ib: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
			args: args{
				data: 77,
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					table: intypes.Table{Name: "table1"},
					errs: intypes.ErrorSlice{
						fmt.Errorf(
							"failed to process argument 0 of Data function: %w",
							fmt.Errorf("recieved value is not a struct type"),
						),
					},
				},
			},
		},
		{
			name: "Error; Bad Additional Argument",
			ib:   insertBuilder{},
			args: args{
				data:     struct{ Data string }{"hi"},
				moreData: []any{"bad_val"},
			},
			want: returningBuilder{
				prevBuilder: insertBuilder{
					errs: intypes.ErrorSlice{
						fmt.Errorf(
							"failed to process argument 1 of Data function: %w",
							fmt.Errorf("recieved value is not a struct type"),
						),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ib.Data(tt.args.data, tt.args.moreData...))
		})
	}
}

func TestNewInsertBuilder(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name string
		args args
		want builders.InsertBuilder
	}{
		{
			name: "Success",
			args: args{
				table: "table1",
			},
			want: insertBuilder{
				table: intypes.Table{Name: "table1"},
			},
		},
		{
			name: "Error; Bad Table Name",
			args: args{
				table: ".bad_name",
			},
			want: insertBuilder{
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse table data from %q: %w", ".bad_name", intypes.ErrMissingSchemaName),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewInsertBuilder(tt.args.table))
		})
	}
}

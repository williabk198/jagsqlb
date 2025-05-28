package inbuilders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	conds "github.com/williabk198/jagsqlb/conditions"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_selectBuilder_Build(t *testing.T) {
	table1 := intypes.Table{Name: "testTable"}
	table1WithSchema := intypes.Table{Name: "testTable", Schema: "testing"}
	table1WithAlias := intypes.Table{Alias: "tt", Name: "testTable"}
	table1WithAliasAndSchema := intypes.Table{Alias: "tt", Name: "testTable", Schema: "testing"}
	column1T1 := intypes.SelectColumn{
		Column: intypes.Column{
			Name:  "testCol1",
			Table: &table1,
		},
	}
	column1T1WithAlias := intypes.SelectColumn{
		Alias: "t1c1",
		Column: intypes.Column{
			Name:  "testCol1",
			Table: &table1WithAlias,
		},
	}
	column1T1All := intypes.SelectColumn{
		Column: intypes.Column{
			Name:  "*",
			Table: &table1,
		},
	}
	column1T1AllWithTableAlias := intypes.SelectColumn{
		Column: intypes.Column{
			Name:  "*",
			Table: &table1WithAlias,
		},
	}

	table2 := intypes.Table{Name: "other"}
	table2WithSchema := intypes.Table{Name: "other", Schema: "public"}
	table2WithAliasAndSchema := intypes.Table{Alias: "o", Name: "other", Schema: "public"}
	column1T2 := intypes.SelectColumn{
		Column: intypes.Column{
			Name:  "testCol1",
			Table: &table2WithAliasAndSchema,
		},
	}
	column1T2WithAlias := intypes.SelectColumn{
		Alias: "t2c1",
		Column: intypes.Column{
			Name:  "testCol1",
			Table: &table2WithSchema,
		},
	}

	tests := []struct {
		name       string
		s          builders.SelectBuilder
		wantQuery  string
		wantParams []any
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name: "Success, 1 Table w/o Columns, Schema & Alias",
			s: selectBuilder{
				tables: []intypes.Table{table1},
			},
			wantQuery: `SELECT FROM "testTable";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, 1 Table w/ Columns; w/o Schema & Alias",
			s: selectBuilder{
				tables:  []intypes.Table{table1},
				columns: []intypes.SelectColumn{column1T1},
			},
			wantQuery: `SELECT "testCol1" FROM "testTable";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, 1 Table w/ Columns & Schema; w/o Alias",
			s: selectBuilder{
				tables:  []intypes.Table{table1WithSchema},
				columns: []intypes.SelectColumn{column1T1},
			},
			wantQuery: `SELECT "testCol1" FROM "testing"."testTable";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, 1 Table w/ Columns, Schema & Alias",
			s: selectBuilder{
				tables:  []intypes.Table{table1WithAliasAndSchema},
				columns: []intypes.SelectColumn{column1T1WithAlias},
			},
			wantQuery: `SELECT "testCol1" AS "t1c1" FROM "testing"."testTable" AS "tt";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, 1 Table w/ All Selector",
			s: selectBuilder{
				tables:  []intypes.Table{table1},
				columns: []intypes.SelectColumn{column1T1All},
			},
			wantQuery: `SELECT * FROM "testTable";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, Multiple Tables w/o Columns, Schemas & Aliases",
			s: selectBuilder{
				tables: []intypes.Table{table1, table2},
			},
			wantQuery: `SELECT FROM "testTable", "other";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, Multiple Tables w/ Columns, Mixed Aliasing & Mixed Schemas",
			s: selectBuilder{
				tables:  []intypes.Table{table1WithAlias, table2WithSchema},
				columns: []intypes.SelectColumn{column1T1AllWithTableAlias, column1T2WithAlias},
			},
			wantQuery: `SELECT "tt".*, "public"."other"."testCol1" AS "t2c1" FROM "testTable" AS "tt", "public"."other";`,
			assertion: assert.NoError,
		},
		{
			name: "Success, Multiple Tables w/ Columns, Mixed Aliasing & Mixed Schemas 2",
			s: selectBuilder{
				tables:  []intypes.Table{table1, table2WithAliasAndSchema},
				columns: []intypes.SelectColumn{column1T1All, column1T2},
			},
			wantQuery: `SELECT "testTable".*, "o"."testCol1" FROM "testTable", "public"."other" AS "o";`,
			assertion: assert.NoError,
		},
		// NOTE: Not going to try to test for every possible error here. That feels like it would be re-testing the parsers.
		//       Instead, just test to see if an error for parseing table, column data, and then both.
		{
			name:      "Error, Bad Table Value",
			s:         NewSelectBuilder(".badValue"),
			assertion: assert.Error,
		},
		{
			name:      "Error, Bad Column Value",
			s:         NewSelectBuilder("testTable", "col1 AS"),
			assertion: assert.Error,
		},
		{
			name:      "Error, Bad Table and Column Value",
			s:         NewSelectBuilder(".testTable", "col1 AS"),
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params, err := tt.s.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wantQuery, query)
			assert.Equal(t, tt.wantParams, params)
		})
	}
}

func Test_selectBuilder_Table(t *testing.T) {
	type args struct {
		table   string
		columns []string
	}

	initialTable := intypes.Table{Alias: "it", Name: "initTable"}
	initialTables := []intypes.Table{initialTable}
	initialColumn := intypes.SelectColumn{Alias: "ic", Column: intypes.Column{Name: "initColumn"}}
	initialColumns := []intypes.SelectColumn{initialColumn}

	testTable := intypes.Table{Name: "testTable"}
	testTableWithSchema := intypes.Table{Name: "testTable", Schema: "testing"}
	testTableWithAlias := intypes.Table{Alias: "tt", Name: "testTable"}
	testTableWithAliasAndSchema := intypes.Table{Alias: "tt", Name: "testTable", Schema: "testing"}

	tests := []struct {
		name string
		s    selectBuilder
		args args
		want builders.SelectBuilder
	}{
		{
			name: "Success w/o Column, Schema & Alias",
			s:    selectBuilder{tables: initialTables, columns: initialColumns},
			args: args{
				table: "testTable",
			},
			want: selectBuilder{
				tables:  []intypes.Table{initialTable, testTable},
				columns: initialColumns,
			},
		},
		{
			name: "Success w/ Column; w/o Schema & Aliases",
			s:    selectBuilder{tables: initialTables, columns: initialColumns},
			args: args{
				table:   "testTable",
				columns: []string{"testColumn"},
			},
			want: selectBuilder{
				tables: []intypes.Table{initialTable, testTable},
				columns: []intypes.SelectColumn{
					initialColumn,
					{
						Column: intypes.Column{Table: &testTable, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column & Aliases; w/o Schema",
			s:    selectBuilder{tables: initialTables, columns: initialColumns},
			args: args{
				table:   "testTable AS tt",
				columns: []string{"testColumn AS tc"},
			},
			want: selectBuilder{
				tables: []intypes.Table{initialTable, testTableWithAlias},
				columns: []intypes.SelectColumn{
					initialColumn,
					{
						Alias:  "tc",
						Column: intypes.Column{Table: &testTableWithAlias, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column, Aliases & Schema",
			s:    selectBuilder{tables: initialTables, columns: initialColumns},
			args: args{
				table:   "testing.testTable AS tt",
				columns: []string{"testColumn AS tc"},
			},
			want: selectBuilder{
				tables: []intypes.Table{initialTable, testTableWithAliasAndSchema},
				columns: []intypes.SelectColumn{
					initialColumn,
					{
						Alias:  "tc",
						Column: intypes.Column{Table: &testTableWithAliasAndSchema, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column & Schema; w/o Aliases",
			s:    selectBuilder{tables: initialTables, columns: initialColumns},
			args: args{
				table:   "testing.testTable",
				columns: []string{"testColumn"},
			},
			want: selectBuilder{
				tables: []intypes.Table{initialTable, testTableWithSchema},
				columns: []intypes.SelectColumn{
					initialColumn,
					{
						Column: intypes.Column{Table: &testTableWithSchema, Name: "testColumn"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.Table(tt.args.table, tt.args.columns...))
		})
	}
}

func Test_selectBuilder_Where(t *testing.T) {
	type args struct {
		cond            inconds.Condition
		additionalConds []inconds.Condition
	}
	testTable1 := intypes.Table{Name: "table1"}
	testSelectCol1 := intypes.SelectColumn{Column: intypes.Column{Table: &testTable1, Name: "col1"}}
	testSelectBuilder := selectBuilder{
		tables:  []intypes.Table{testTable1},
		columns: []intypes.SelectColumn{testSelectCol1},
	}
	cond1 := conds.Equals("col1", "testing")
	cond2 := conds.Between("col2", 42, 56)
	cond3 := conds.GraterThan("col3", 98.76)

	tests := []struct {
		name string
		s    selectBuilder
		args args
		want builders.WhereBuilder
	}{
		{
			name: "Success; Minimal",
			s:    testSelectBuilder,
			args: args{
				cond: cond1,
			},
			want: whereBuilder{
				mainQuery: testSelectBuilder,
				conditions: []whereCondition{
					{
						condition: cond1,
					},
				},
			},
		},
		{
			name: "Success; Multiple Conditions",
			s:    testSelectBuilder,
			args: args{
				cond:            cond1,
				additionalConds: []inconds.Condition{cond2, cond3},
			},
			want: whereBuilder{
				mainQuery: testSelectBuilder,
				conditions: []whereCondition{
					{condition: cond1},
					{condition: cond2, conjunction: "AND"},
					{condition: cond3, conjunction: "AND"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.Where(tt.args.cond, tt.args.additionalConds...))
		})
	}
}

func TestNewSelectBuilder(t *testing.T) {
	type args struct {
		table   string
		columns []string
	}

	testTable := intypes.Table{Name: "testTable"}
	testTableWithSchema := intypes.Table{Name: "testTable", Schema: "testing"}
	testTableWithAlias := intypes.Table{Alias: "tt", Name: "testTable"}
	testTableWithAliasAndSchema := intypes.Table{Alias: "tt", Name: "testTable", Schema: "testing"}

	tests := []struct {
		name string
		args args
		want builders.SelectBuilder
	}{
		{
			name: "Success w/o Column, Schema & Alias",
			args: args{
				table: "testTable",
			},
			want: selectBuilder{
				tables: []intypes.Table{testTable},
			},
		},
		{
			name: "Success w/ Column; w/o Schema & Aliases",
			args: args{
				table:   "testTable",
				columns: []string{"testColumn"},
			},
			want: selectBuilder{
				tables: []intypes.Table{testTable},
				columns: []intypes.SelectColumn{
					{
						Column: intypes.Column{Table: &testTable, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column & Aliases; w/o Schema",
			args: args{
				table:   "testTable AS tt",
				columns: []string{"testColumn AS tc"},
			},
			want: selectBuilder{
				tables: []intypes.Table{testTableWithAlias},
				columns: []intypes.SelectColumn{
					{
						Alias:  "tc",
						Column: intypes.Column{Table: &testTableWithAlias, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column, Aliases & Schema",
			args: args{
				table:   "testing.testTable AS tt",
				columns: []string{"testColumn AS tc"},
			},
			want: selectBuilder{
				tables: []intypes.Table{testTableWithAliasAndSchema},
				columns: []intypes.SelectColumn{
					{
						Alias:  "tc",
						Column: intypes.Column{Table: &testTableWithAliasAndSchema, Name: "testColumn"},
					},
				},
			},
		},
		{
			name: "Success w/ Column & Schema; w/o Aliases",
			args: args{
				table:   "testing.testTable",
				columns: []string{"testColumn"},
			},
			want: selectBuilder{
				tables: []intypes.Table{testTableWithSchema},
				columns: []intypes.SelectColumn{
					{
						Column: intypes.Column{Table: &testTableWithSchema, Name: "testColumn"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSelectBuilder(tt.args.table, tt.args.columns...)
			assert.Equal(t, tt.want, got)
		})
	}
}

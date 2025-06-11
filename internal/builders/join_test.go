package inbuilders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	conds "github.com/williabk198/jagsqlb/conditions"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	"github.com/williabk198/jagsqlb/join"
	"github.com/williabk198/jagsqlb/types"
)

func Test_joinBuilder_Build(t *testing.T) {

	type wants struct {
		query       string
		queryParams []any
	}

	testTable1 := intypes.Table{
		Alias: "t1",
		Name:  "table1",
	}
	testTable2 := intypes.Table{
		Alias: "t2",
		Name:  "table2",
	}
	testTable3 := intypes.Table{
		Alias: "t3",
		Name:  "table3",
	}

	testSelectColumn1 := intypes.SelectColumn{Column: intypes.Column{Name: "*", Table: &testTable1}}
	testSelectColumn2 := intypes.SelectColumn{Alias: "t2c1", Column: intypes.Column{Name: "col1", Table: &testTable2}}

	testJoinCondition1 := joinCondition{
		joinTable: testTable2,
		joinType:  join.TypeInner,
		joinRelation: injoin.Relation{
			Keyword:  "USING",
			Relation: "col1",
		},
	}
	testJoinCondition2 := joinCondition{
		joinTable: testTable3,
		joinType:  join.TypeLeft,
		joinRelation: injoin.Relation{
			Keyword: "ON",
			Relation: []inconds.Condition{
				conds.Equals("t1.col1", conds.ColumnValue("t3.col2")),
			},
		},
	}
	testJoinCondition3 := joinCondition{
		joinTable: testTable2,
		joinType:  join.TypeRight,
		joinRelation: injoin.Relation{
			Keyword:  "USING",
			Relation: "col2",
		},
	}

	tests := []struct {
		name      string
		jb        joinBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Single Join",
			jb: joinBuilder{
				selectBuilder: selectBuilder{
					tables:  []intypes.Table{testTable1},
					columns: []intypes.SelectColumn{testSelectColumn1, testSelectColumn2},
				},
				joins: []joinCondition{testJoinCondition1},
			},
			wants: wants{
				query: `SELECT "t1".*, "t2"."col1" AS "t2c1" FROM "table1" AS "t1" INNER JOIN "table2" AS "t2" USING ("col1");`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Multiple Join",
			jb: joinBuilder{
				selectBuilder: selectBuilder{
					tables:  []intypes.Table{testTable1},
					columns: []intypes.SelectColumn{testSelectColumn1},
				},
				joins: []joinCondition{testJoinCondition2, testJoinCondition3},
			},
			wants: wants{
				query: `SELECT "t1".* FROM "table1" AS "t1" LEFT JOIN "table3" AS "t3" ON "t1"."col1" = "t3"."col2" RIGHT JOIN "table2" AS "t2" USING ("col2");`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Error Slice not Empty",
			jb: joinBuilder{
				errs: intypes.ErrorSlice{assert.AnError},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; USING w/ Invalid Column Definition",
			jb: joinBuilder{
				selectBuilder: selectBuilder{},
				joins: []joinCondition{
					{
						joinTable:    testTable1,
						joinType:     join.TypeCross,
						joinRelation: join.Using(".bad_col"),
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; ON w/ Invalid Condition",
			jb: joinBuilder{
				joins: []joinCondition{
					{
						joinTable:    testTable1,
						joinType:     join.TypeOutter,
						joinRelation: join.On(conds.Equals(".bad_col", 42)),
					},
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; ON w/ Invalid Condition 2",
			jb: joinBuilder{
				joins: []joinCondition{
					{
						joinTable: testTable1,
						joinType:  join.TypeRight,
						joinRelation: join.On(
							conds.Equals("col1", conds.ColumnValue("t2.col2")),
							conds.GreaterThan(".col2", 55),
							conds.Equals(".bad_col", 42),
						),
					},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotQueryParams, err := tt.jb.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.queryParams, gotQueryParams)
		})
	}
}

func Test_joinBuilder_Join(t *testing.T) {
	type args struct {
		joinType       injoin.Type
		table          string
		joinRelation   injoin.Relation
		includeColumns []string
	}

	testTable1 := intypes.Table{
		Alias: "t1",
		Name:  "table1",
	}
	testTable2 := intypes.Table{
		Alias: "t2",
		Name:  "table2",
	}
	testTable3 := intypes.Table{
		Alias: "t3",
		Name:  "table3",
	}

	testSelectColumn1 := intypes.SelectColumn{Column: intypes.Column{Name: "*", Table: &testTable1}}

	testJoinCondition1 := joinCondition{
		joinTable: testTable2,
		joinType:  join.TypeInner,
		joinRelation: injoin.Relation{
			Keyword:  "USING",
			Relation: "col1",
		},
	}
	testJoinCondition2 := joinCondition{
		joinTable: testTable3,
		joinType:  join.TypeLeft,
		joinRelation: injoin.Relation{
			Keyword: "ON",
			Relation: []inconds.Condition{
				conds.Equals("t1.col1", conds.ColumnValue("t3.col2")),
			},
		},
	}
	testJoinCondition3 := joinCondition{
		joinTable: testTable2,
		joinType:  join.TypeRight,
		joinRelation: injoin.Relation{
			Keyword:  "USING",
			Relation: "col2",
		},
	}

	tests := []struct {
		name string
		jb   joinBuilder
		args args
		want builders.JoinBuilder
	}{
		{
			name: "Success; No Additional Columns",
			jb: joinBuilder{
				selectBuilder: selectBuilder{},
				joins:         []joinCondition{testJoinCondition1},
			},
			args: args{
				joinType:     join.TypeLeft,
				table:        "table3 AS t3",
				joinRelation: join.On(conds.Equals("t1.col1", conds.ColumnValue("t3.col2"))),
			},
			want: joinBuilder{
				selectBuilder: selectBuilder{},
				joins:         []joinCondition{testJoinCondition1, testJoinCondition2},
			},
		},
		{
			name: "Success; With Additional Columns",
			jb: joinBuilder{
				selectBuilder: selectBuilder{
					tables:  []intypes.Table{testTable1},
					columns: []intypes.SelectColumn{testSelectColumn1},
				},
				joins: []joinCondition{testJoinCondition1},
			},
			args: args{
				joinType:       join.TypeRight,
				table:          "table2 AS t2",
				joinRelation:   join.Using("col2"),
				includeColumns: []string{"col3"},
			},
			want: joinBuilder{
				selectBuilder: selectBuilder{
					tables: []intypes.Table{testTable1},
					columns: []intypes.SelectColumn{
						testSelectColumn1,
						{Column: intypes.Column{Name: "col3", Table: &testTable2}},
					},
				},
				joins: []joinCondition{testJoinCondition1, testJoinCondition3},
			},
		},
		{
			name: "Error; Bad Table Name",
			jb: joinBuilder{
				selectBuilder: selectBuilder{},
				joins:         []joinCondition{},
			},
			args: args{
				joinType:     join.TypeLeft,
				table:        ".bad_table",
				joinRelation: join.Using("col1"),
			},
			want: joinBuilder{
				selectBuilder: selectBuilder{},
				joins:         []joinCondition{},
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse table in JOIN clause: %w", fmt.Errorf("failed to parse table data from %q: %w", ".bad_table", intypes.ErrMissingSchemaName)),
				},
			},
		},
		{
			name: "Error; Bad Additional Column Name",
			jb: joinBuilder{
				selectBuilder: selectBuilder{},
				joins: []joinCondition{
					testJoinCondition1,
				},
			},
			args: args{
				joinType: join.TypeInner,
				table:    "table2 AS t2",
				joinRelation: injoin.Relation{
					Keyword:  "USING",
					Relation: "col1",
				},
				includeColumns: []string{".bad_col"},
			},
			want: joinBuilder{
				selectBuilder: selectBuilder{},
				joins: []joinCondition{
					testJoinCondition1,
				},
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse column %q in %s of %s: %w", ".bad_col", join.TypeInner, "table2", fmt.Errorf("failed to parse table data provided in %q: %w", ".bad_col", intypes.ErrMissingTableName)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.jb.Join(tt.args.joinType, tt.args.table, tt.args.joinRelation, tt.args.includeColumns...))
		})
	}
}

func Test_joinBuilder_Where(t *testing.T) {
	type args struct {
		condition      inconds.Condition
		moreConditions []inconds.Condition
	}
	tests := []struct {
		name string
		jb   joinBuilder
		args args
		want builders.WhereBuilder
	}{
		{
			name: "Success",
			jb:   joinBuilder{},
			args: args{
				condition: conds.GreaterThanEqual("t2.col3", 59),
			},
			want: whereBuilder{
				mainQuery: joinBuilder{},
				conditions: []whereCondition{
					{condition: conds.GreaterThanEqual("t2.col3", 59)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.jb.Where(tt.args.condition, tt.args.moreConditions...))
		})
	}
}

func Test_joinBuilder_Limit(t *testing.T) {
	type args struct {
		limit uint
	}
	tests := []struct {
		name string
		jb   joinBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			jb:   joinBuilder{},
			args: args{
				limit: 20,
			},
			want: limitBuilder{
				precedingBuilder: joinBuilder{},
				limit:            20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.jb.Limit(tt.args.limit))
		})
	}
}

func Test_joinBuilder_Offset(t *testing.T) {
	type args struct {
		offset uint
	}
	tests := []struct {
		name string
		jb   joinBuilder
		args args
		want builders.OffsetBuilder
	}{
		{
			name: "Success",
			jb:   joinBuilder{},
			args: args{
				offset: 100,
			},
			want: offsetBuilder{
				precedingBuilder: joinBuilder{},
				offset:           100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.jb.Offset(tt.args.offset))
		})
	}
}

func Test_joinBuilder_OrderBy(t *testing.T) {
	type args struct {
		ordering      types.ColumnOrdering
		moreOrderings []types.ColumnOrdering
	}
	tests := []struct {
		name string
		jb   joinBuilder
		args args
		want builders.OrderByBuilder
	}{
		{
			name: "Success",
			jb:   joinBuilder{},
			args: args{
				ordering: types.ColumnOrdering{ColumnName: "t1.col1", Ordering: types.OrderingDescending},
			},
			want: orderByBuilder{
				precedingBuilder: joinBuilder{},
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: "t1.col1", Ordering: types.OrderingDescending},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.jb.OrderBy(tt.args.ordering, tt.args.moreOrderings...))
		})
	}
}

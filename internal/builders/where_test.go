package inbuilders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	"github.com/williabk198/jagsqlb/condition"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	intypes "github.com/williabk198/jagsqlb/internal/types"
	"github.com/williabk198/jagsqlb/types"
)

func Test_whereBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		w         selectWhereBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Simple Conditions",
			w: selectWhereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1", "col1", "col2"),
				conditions: []whereCondition{
					{condition: condition.Equals("col1", "test")},
					{condition: condition.NotBetween("col2", 10, 23), conjunction: "OR"},
				},
			},
			wants: wants{
				query:  `SELECT "col1", "col2" FROM "table1" AS "t1" WHERE "col1" = $1 OR "col2" NOT BETWEEN $2 AND $3;`,
				params: []any{"test", 10, 23},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Simple Conditions w/ ColumnValue",
			w: selectWhereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1").Table("table2 AS t2", "col1"),
				conditions: []whereCondition{
					{condition: condition.Equals("t1.col1", incondition.ColumnValue{ColumnName: "t2.col2"})},
				},
			},
			wants: wants{
				query:  `SELECT "t2"."col1" FROM "table1" AS "t1", "table2" AS "t2" WHERE "t1"."col1" = "t2"."col2";`,
				params: nil,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Grouped Conditions",
			w: selectWhereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: condition.GroupedOr(condition.Equals("col1", "test"), condition.GreaterThanEqual("col2", 52))},
					{condition: condition.GroupedOr(condition.NotIn("col3", []any{"test", "testing"}), condition.LessThan("col2", 52)), conjunction: "AND"},
				},
			},
			wants: wants{
				query:  `SELECT * FROM "table1" WHERE ("col1" = $1 OR "col2" >= $2) AND ("col3" NOT IN $3 OR "col2" < $4);`,
				params: []any{"test", 52, []any{"test", "testing"}, 52},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Mixed Conditions",
			w: selectWhereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: condition.GroupedAnd(condition.Equals("col1", "test"), condition.GreaterThanEqual("col2", 52))},
					{condition: condition.LessThan("col3", 128), conjunction: "OR"},
				},
			},
			wants: wants{
				query:  `SELECT * FROM "table1" WHERE ("col1" = $1 AND "col2" >= $2) OR "col3" < $3;`,
				params: []any{"test", 52, 128},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotQueryParams, err := tt.w.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotQueryParams)
		})
	}
}

func Test_whereBuilder_And(t *testing.T) {
	type args struct {
		cond            incondition.Condition
		additionalConds []incondition.Condition
	}

	testWhereCond1 := whereCondition{
		condition: condition.Equals("t1.col1", "testing"),
	}
	testCondInput1 := condition.LessThan("col2", 98.76)
	testCondInput2 := condition.Between("col3", 17, 78)

	tests := []struct {
		name string
		w    selectWhereBuilder
		args args
		want builders.SelectWhereBuilder
	}{
		{
			name: "Single Condition",
			w: selectWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: selectWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "AND", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			w: selectWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []incondition.Condition{testCondInput2},
			},
			want: selectWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "AND", condition: testCondInput1},
					{conjunction: "AND", condition: testCondInput2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.w.And(tt.args.cond, tt.args.additionalConds...))
		})
	}
}

func Test_whereBuilder_Or(t *testing.T) {
	type args struct {
		cond            incondition.Condition
		additionalConds []incondition.Condition
	}

	testWhereCond1 := whereCondition{
		condition: condition.Equals("t1.col1", "testing"),
	}
	testCondInput1 := condition.LessThan("col2", 98.76)
	testCondInput2 := condition.Between("col3", 17, 78)

	tests := []struct {
		name string
		w    selectWhereBuilder
		args args
		want builders.SelectWhereBuilder
	}{
		{
			name: "Single Condition",
			w: selectWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: selectWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "OR", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			w: selectWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []incondition.Condition{testCondInput2},
			},
			want: selectWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "OR", condition: testCondInput1},
					{conjunction: "OR", condition: testCondInput2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.w.Or(tt.args.cond, tt.args.additionalConds...))
		})
	}
}

func Test_whereBuilder_Limit(t *testing.T) {
	type args struct {
		limit uint
	}
	tests := []struct {
		name string
		w    selectWhereBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			w:    selectWhereBuilder{},
			args: args{
				limit: 100,
			},
			want: limitBuilder{
				precedingBuilder: selectWhereBuilder{},
				limit:            100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.w.Limit(tt.args.limit))
		})
	}
}

func Test_whereBuilder_Offset(t *testing.T) {
	type args struct {
		offset uint
	}
	tests := []struct {
		name string
		w    selectWhereBuilder
		args args
		want builders.OffsetBuilder
	}{
		{
			name: "Success",
			w:    selectWhereBuilder{},
			args: args{
				offset: 100,
			},
			want: offsetBuilder{
				precedingBuilder: selectWhereBuilder{},
				offset:           100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.w.Offset(tt.args.offset))
		})
	}
}

func Test_whereBuilder_OrderBy(t *testing.T) {
	type args struct {
		ordering      types.ColumnOrdering
		moreOrderings []types.ColumnOrdering
	}
	tests := []struct {
		name string
		w    selectWhereBuilder
		args args
		want builders.OrderByBuilder
	}{
		{
			name: "Success",
			w:    selectWhereBuilder{},
			args: args{
				ordering: types.ColumnOrdering{ColumnName: "col1", Ordering: types.OrderingAscending},
			},
			want: orderByBuilder{
				precedingBuilder: selectWhereBuilder{},
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: "col1", Ordering: types.OrderingAscending},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.w.OrderBy(tt.args.ordering, tt.args.moreOrderings...))
		})
	}
}

func Test_returningWhereBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		rwb       returningWhereBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Simple Conditions",
			rwb: returningWhereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1", "col1", "col2"),
				conditions: []whereCondition{
					{condition: condition.Equals("col1", "test")},
					{condition: condition.NotBetween("col2", 10, 23), conjunction: "OR"},
				},
			},
			wants: wants{
				query:  `SELECT "col1", "col2" FROM "table1" AS "t1" WHERE "col1" = $1 OR "col2" NOT BETWEEN $2 AND $3;`,
				params: []any{"test", 10, 23},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Simple Conditions w/ ColumnValue",
			rwb: returningWhereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1").Table("table2 AS t2", "col1"),
				conditions: []whereCondition{
					{condition: condition.Equals("t1.col1", incondition.ColumnValue{ColumnName: "t2.col2"})},
				},
			},
			wants: wants{
				query:  `SELECT "t2"."col1" FROM "table1" AS "t1", "table2" AS "t2" WHERE "t1"."col1" = "t2"."col2";`,
				params: nil,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Grouped Conditions",
			rwb: returningWhereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: condition.GroupedOr(condition.Equals("col1", "test"), condition.GreaterThanEqual("col2", 52))},
					{condition: condition.GroupedOr(condition.NotIn("col3", []any{"test", "testing"}), condition.LessThan("col2", 52)), conjunction: "AND"},
				},
			},
			wants: wants{
				query:  `SELECT * FROM "table1" WHERE ("col1" = $1 OR "col2" >= $2) AND ("col3" NOT IN $3 OR "col2" < $4);`,
				params: []any{"test", 52, []any{"test", "testing"}, 52},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Mixed Conditions",
			rwb: returningWhereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: condition.GroupedAnd(condition.Equals("col1", "test"), condition.GreaterThanEqual("col2", 52))},
					{condition: condition.LessThan("col3", 128), conjunction: "OR"},
				},
			},
			wants: wants{
				query:  `SELECT * FROM "table1" WHERE ("col1" = $1 AND "col2" >= $2) OR "col3" < $3;`,
				params: []any{"test", 52, 128},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotQueryParams, err := tt.rwb.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotQueryParams)
		})
	}
}

func Test_returningWhereBuilder_And(t *testing.T) {
	type args struct {
		cond            incondition.Condition
		additionalConds []incondition.Condition
	}

	testWhereCond1 := whereCondition{
		condition: condition.Equals("t1.col1", "testing"),
	}
	testCondInput1 := condition.LessThan("col2", 98.76)
	testCondInput2 := condition.Between("col3", 17, 78)

	tests := []struct {
		name string
		rwb  returningWhereBuilder
		args args
		want builders.ReturningWhereBuilder
	}{
		{
			name: "Single Condition",
			rwb: returningWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: returningWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "AND", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			rwb: returningWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []incondition.Condition{testCondInput2},
			},
			want: returningWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "AND", condition: testCondInput1},
					{conjunction: "AND", condition: testCondInput2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rwb.And(tt.args.cond, tt.args.additionalConds...))
		})
	}
}

func Test_returningWhereBuilder_Or(t *testing.T) {
	type args struct {
		cond            incondition.Condition
		additionalConds []incondition.Condition
	}

	testWhereCond1 := whereCondition{
		condition: condition.Equals("t1.col1", "testing"),
	}
	testCondInput1 := condition.LessThan("col2", 98.76)
	testCondInput2 := condition.Between("col3", 17, 78)

	tests := []struct {
		name string
		rwb  returningWhereBuilder
		args args
		want builders.ReturningWhereBuilder
	}{
		{
			name: "Single Condition",
			rwb: returningWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: returningWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "OR", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			rwb: returningWhereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []incondition.Condition{testCondInput2},
			},
			want: returningWhereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "OR", condition: testCondInput1},
					{conjunction: "OR", condition: testCondInput2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rwb.Or(tt.args.cond, tt.args.additionalConds...))
		})
	}
}

func Test_returningWhereBuilder_Returning(t *testing.T) {
	type args struct {
		column      string
		moreColumns []string
	}

	testWhereCond := returningWhereBuilder{
		conditions: []whereCondition{},
	}

	tests := []struct {
		name string
		rwb  returningWhereBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Single Column",
			rwb:  testWhereCond,
			args: args{
				column: "*",
			},
			want: returningBuilder{
				prevBuilder: testWhereCond,
				returningColumns: []intypes.Column{
					{Name: "*"},
				},
			},
		},
		{
			name: "Multiple Columns",
			rwb:  testWhereCond,
			args: args{
				column: "col1",
				moreColumns: []string{
					"col2", "col3",
				},
			},
			want: returningBuilder{
				prevBuilder: testWhereCond,
				returningColumns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
					{Name: "col3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rwb.Returning(tt.args.column, tt.args.moreColumns...))
		})
	}
}

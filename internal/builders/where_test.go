package inbuilders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	conds "github.com/williabk198/jagsqlb/conditions"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	"github.com/williabk198/jagsqlb/types"
)

func Test_whereBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		w         whereBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Simple Conditions",
			w: whereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1", "col1", "col2"),
				conditions: []whereCondition{
					{condition: conds.Equals("col1", "test")},
					{condition: conds.NotBetween("col2", 10, 23), conjunction: "OR"},
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
			w: whereBuilder{
				mainQuery: NewSelectBuilder("table1 AS t1").Table("table2 AS t2", "col1"),
				conditions: []whereCondition{
					{condition: conds.Equals("t1.col1", inconds.ColumnValue{ColumnName: "t2.col2"})},
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
			w: whereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: conds.GroupedOr(conds.Equals("col1", "test"), conds.GreaterThanEqual("col2", 52))},
					{condition: conds.GroupedOr(conds.NotIn("col3", []any{"test", "testing"}), conds.LessThan("col2", 52)), conjunction: "AND"},
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
			w: whereBuilder{
				mainQuery: NewSelectBuilder("table1", "*"),
				conditions: []whereCondition{
					{condition: conds.GroupedAnd(conds.Equals("col1", "test"), conds.GreaterThanEqual("col2", 52))},
					{condition: conds.LessThan("col3", 128), conjunction: "OR"},
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
		cond            inconds.Condition
		additionalConds []inconds.Condition
	}

	testWhereCond1 := whereCondition{
		condition: conds.Equals("t1.col1", "testing"),
	}
	testCondInput1 := conds.LessThan("col2", 98.76)
	testCondInput2 := conds.Between("col3", 17, 78)

	tests := []struct {
		name string
		w    whereBuilder
		args args
		want builders.WhereBuilder
	}{
		{
			name: "Single Condition",
			w: whereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: whereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "AND", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			w: whereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []inconds.Condition{testCondInput2},
			},
			want: whereBuilder{
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
		cond            inconds.Condition
		additionalConds []inconds.Condition
	}

	testWhereCond1 := whereCondition{
		condition: conds.Equals("t1.col1", "testing"),
	}
	testCondInput1 := conds.LessThan("col2", 98.76)
	testCondInput2 := conds.Between("col3", 17, 78)

	tests := []struct {
		name string
		w    whereBuilder
		args args
		want builders.WhereBuilder
	}{
		{
			name: "Single Condition",
			w: whereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond: testCondInput1,
			},
			want: whereBuilder{
				conditions: []whereCondition{
					testWhereCond1,
					{conjunction: "OR", condition: testCondInput1},
				},
			},
		},
		{
			name: "Multiple Conditions",
			w: whereBuilder{
				conditions: []whereCondition{testWhereCond1},
			},
			args: args{
				cond:            testCondInput1,
				additionalConds: []inconds.Condition{testCondInput2},
			},
			want: whereBuilder{
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
		w    whereBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			w:    whereBuilder{},
			args: args{
				limit: 100,
			},
			want: limitBuilder{
				precedingBuilder: whereBuilder{},
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
		w    whereBuilder
		args args
		want builders.OffsetBuilder
	}{
		{
			name: "Success",
			w:    whereBuilder{},
			args: args{
				offset: 100,
			},
			want: offsetBuilder{
				precedingBuilder: whereBuilder{},
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
		w    whereBuilder
		args args
		want builders.OrderByBuilder
	}{
		{
			name: "Success",
			w:    whereBuilder{},
			args: args{
				ordering: types.ColumnOrdering{ColumnName: "col1", Ordering: types.OrderingAscending},
			},
			want: orderByBuilder{
				precedingBuilder: whereBuilder{},
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

package inbuilders

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	"github.com/williabk198/jagsqlb/types"
)

func Test_orderByBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		obb       orderByBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Single Asecending",
			obb: orderByBuilder{
				precedingBuilder: NewSelectBuilder("table1", "*"),
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: "column1", Ordering: types.OrderingAscending},
				},
			},
			wants: wants{
				query: `SELECT * FROM "table1" ORDER BY "column1" ASC;`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Single Descending",
			obb: orderByBuilder{
				precedingBuilder: NewSelectBuilder("table1", "*"),
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: "column1", Ordering: types.OrderingDescending},
				},
			},
			wants: wants{
				query: `SELECT * FROM "table1" ORDER BY "column1" DESC;`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Multiple Mixed",
			obb: orderByBuilder{
				precedingBuilder: NewSelectBuilder("table1", "*"),
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: "column1", Ordering: types.OrderingAscending},
					{ColumnName: "column2", Ordering: types.OrderingDescending},
				},
			},
			wants: wants{
				query: `SELECT * FROM "table1" ORDER BY "column1" ASC, "column2" DESC;`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Preceding Builder",
			obb: orderByBuilder{
				precedingBuilder: NewSelectBuilder(".t1 AS", "col1"),
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Bad Column in Ordering",
			obb: orderByBuilder{
				precedingBuilder: NewSelectBuilder("table1 AS t1", "*"),
				columnOrderings: []types.ColumnOrdering{
					{ColumnName: ".col1", Ordering: types.OrderingDescending},
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.obb.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotParams)
		})
	}
}

func Test_orderByBuilder_Offset(t *testing.T) {
	type args struct {
		offset uint
	}

	testOrderBuilder := orderByBuilder{
		precedingBuilder: NewSelectBuilder("table1", "*"),
		columnOrderings: []types.ColumnOrdering{
			{ColumnName: "col1", Ordering: types.OrderingAscending},
		},
	}

	tests := []struct {
		name string
		oob  orderByBuilder
		args args
		want builders.OffsetBuilder
	}{
		{
			name: "Success",
			oob:  testOrderBuilder,
			args: args{
				offset: 10,
			},
			want: offsetBuilder{
				precedingBuilder: testOrderBuilder,
				offset:           10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.oob.Offset(tt.args.offset))
		})
	}
}

func Test_orderByBuilder_Limit(t *testing.T) {
	type args struct {
		limit uint
	}

	testOrderBuilder := orderByBuilder{
		precedingBuilder: NewSelectBuilder("table1", "*"),
		columnOrderings: []types.ColumnOrdering{
			{ColumnName: "col1", Ordering: types.OrderingAscending},
		},
	}

	tests := []struct {
		name string
		oob  orderByBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			oob:  testOrderBuilder,
			args: args{
				limit: 5,
			},
			want: limitBuilder{
				precedingBuilder: testOrderBuilder,
				limit:            5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.oob.Limit(tt.args.limit))
		})
	}
}

func Test_offsetBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		ob        offsetBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			ob: offsetBuilder{
				precedingBuilder: NewSelectBuilder("table1", "*"),
				offset:           100,
			},
			wants: wants{
				query: `SELECT * FROM "table1" OFFSET 100;`,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.ob.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotParams)
		})
	}
}

func Test_offsetBuilder_Limit(t *testing.T) {
	type args struct {
		limit uint
	}

	testOffsetBuilder := offsetBuilder{
		precedingBuilder: NewSelectBuilder("table1", "*"),
		offset:           50,
	}

	tests := []struct {
		name string
		ob   offsetBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			ob:   testOffsetBuilder,
			args: args{
				limit: 10,
			},
			want: limitBuilder{
				precedingBuilder: testOffsetBuilder,
				limit:            10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ob.Limit(tt.args.limit))
		})
	}
}

func Test_limitBuilder_Build(t *testing.T) {
	type wants struct {
		query string
		parms []any
	}

	tests := []struct {
		name      string
		lb        limitBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			lb: limitBuilder{
				precedingBuilder: NewSelectBuilder("table1", "col1"),
				limit:            25,
			},
			wants: wants{
				query: `SELECT "col1" FROM "table1" LIMIT 25`,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.lb.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.parms, gotParams)
		})
	}
}

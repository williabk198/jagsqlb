package incondition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCondition_Parameterize(t *testing.T) {

	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		sc        SimpleCondition
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Equals",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "=",
				Values:     []any{"test"},
			},
			wants: wants{
				query:  `"col1" = ?`,
				params: []any{"test"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Equals w/ ColumnValue",
			sc: SimpleCondition{
				ColumnName: "t1.col1",
				Operator:   "=",
				Values:     []any{ColumnValue{ColumnName: "t2.col2"}},
			},
			wants: wants{
				query:  `"t1"."col1" = "t2"."col2"`,
				params: []any{},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Not Equals",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "!=",
				Values:     []any{"test"},
			},
			wants: wants{
				query:  `"col1" != ?`,
				params: []any{"test"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Greater Than",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   ">",
				Values:     []any{42},
			},
			wants: wants{
				query:  `"col1" > ?`,
				params: []any{42},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Greater Than Equals",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   ">=",
				Values:     []any{42},
			},
			wants: wants{
				query:  `"col1" >= ?`,
				params: []any{42},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Less Than",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "<",
				Values:     []any{42},
			},
			wants: wants{
				query:  `"col1" < ?`,
				params: []any{42},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Less Than Equal",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "<=",
				Values:     []any{42},
			},
			wants: wants{
				query:  `"col1" <= ?`,
				params: []any{42},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Is Null",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "IS",
				Values:     []any{"NULL"},
			},
			wants: wants{
				query: `"col1" IS NULL`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Is Not Null",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "IS NOT",
				Values:     []any{"NULL"},
			},
			wants: wants{
				query: `"col1" IS NOT NULL`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; In",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "IN",
				Values:     []any{"test", "testing"},
			},
			wants: wants{
				query:  `"col1" IN ?`,
				params: []any{"test", "testing"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Not In",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "NOT IN",
				Values:     []any{"test", "testing"},
			},
			wants: wants{
				query:  `"col1" NOT IN ?`,
				params: []any{"test", "testing"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Between",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "BETWEEN",
				Values:     []any{42, 56},
			},
			wants: wants{
				query:  `"col1" BETWEEN ? AND ?`,
				params: []any{42, 56},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Between with First ColumnValue",
			sc: SimpleCondition{
				ColumnName: "t1.col1",
				Operator:   "BETWEEN",
				Values:     []any{ColumnValue{ColumnName: "t2.col1"}, 83},
			},
			wants: wants{
				query:  `"t1"."col1" BETWEEN "t2"."col1" AND ?`,
				params: []any{83},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Between with Second ColumnValue",
			sc: SimpleCondition{
				ColumnName: "t1.col1",
				Operator:   "BETWEEN",
				Values:     []any{83, ColumnValue{ColumnName: "t2.col1"}},
			},
			wants: wants{
				query:  `"t1"."col1" BETWEEN ? AND "t2"."col1"`,
				params: []any{83},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success: Between with Two ColumnValues",
			sc: SimpleCondition{
				ColumnName: "t1.col1",
				Operator:   "BETWEEN",
				Values:     []any{ColumnValue{ColumnName: "t2.col1"}, ColumnValue{ColumnName: "t2.col2"}},
			},
			wants: wants{
				query:  `"t1"."col1" BETWEEN "t2"."col1" AND "t2"."col2"`,
				params: []any{},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Not Between",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "NOT BETWEEN",
				Values:     []any{42, 56},
			},
			wants: wants{
				query:  `"col1" NOT BETWEEN ? AND ?`,
				params: []any{42, 56},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad Column Definition",
			sc: SimpleCondition{
				ColumnName: ".badColumn",
				Operator:   "=",
				Values:     []any{"n/a"},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; In Condition Contains ColumnValue",
			sc: SimpleCondition{
				ColumnName: "col1",
				Operator:   "NOT IN",
				Values:     []any{52, ColumnValue{ColumnName: "t2.col2"}, 77},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.sc.Parameterize()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotParams)
		})
	}
}

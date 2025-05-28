package inconds

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
			name: "Success; In",
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

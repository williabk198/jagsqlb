package intypes

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
				Operator:   "BETWEEN",
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

func TestGroupedConditions_Parameterize(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	testCond1 := SimpleCondition{
		ColumnName: "t1.col1",
		Operator:   "NOT IN",
		Values:     []any{"test", "testing", "tester"},
	}
	testCond2 := SimpleCondition{
		ColumnName: "col2",
		Operator:   "BETWEEN",
		Values:     []any{42, 56},
	}
	testCond3 := SimpleCondition{
		ColumnName: "col3",
		Operator:   "!=",
		Values:     []any{23},
	}
	testCond4 := SimpleCondition{
		ColumnName: "col4",
		Operator:   ">=",
		Values:     []any{12.34},
	}

	testGroupCond1 := GroupedConditions{
		Conjunction: " OR ",
		Conditions:  []Condition{testCond2, testCond3},
	}
	testGroupCond2 := GroupedConditions{
		Conjunction: " AND ",
		Conditions:  []Condition{testCond1, testCond4},
	}

	tests := []struct {
		name      string
		gc        GroupedConditions
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Grouped AND",
			gc: GroupedConditions{
				Conjunction: " AND ",
				Conditions:  []Condition{testCond1, testCond2},
			},
			wants: wants{
				query:  `("t1"."col1" NOT IN ? AND "col2" BETWEEN ? AND ?)`,
				params: []any{[]any{"test", "testing", "tester"}, 42, 56},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Grouped AND with Nested Group",
			gc: GroupedConditions{
				Conjunction: " AND ",
				Conditions:  []Condition{testGroupCond1, testCond1, testCond4},
			},
			wants: wants{
				query:  `(("col2" BETWEEN ? AND ? OR "col3" != ?) AND "t1"."col1" NOT IN ? AND "col4" >= ?)`,
				params: []any{42, 56, 23, []any{"test", "testing", "tester"}, 12.34},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Grouped OR",
			gc: GroupedConditions{
				Conjunction: " OR ",
				Conditions:  []Condition{testCond3, testCond4},
			},
			wants: wants{
				query:  `("col3" != ? OR "col4" >= ?)`,
				params: []any{23, 12.34},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Grouped OR wih Nested Group",
			gc: GroupedConditions{
				Conjunction: " OR ",
				Conditions:  []Condition{testCond2, testGroupCond2, testCond3},
			},
			wants: wants{
				query:  `("col2" BETWEEN ? AND ? OR ("t1"."col1" NOT IN ? AND "col4" >= ?) OR "col3" != ?)`,
				params: []any{42, 56, []any{"test", "testing", "tester"}, 12.34, 23},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotParams, err := tt.gc.Parameterize()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotParams)
		})
	}
}

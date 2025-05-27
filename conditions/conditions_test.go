package conds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func TestEquals(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      "testing",
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "=",
				Values:     []any{"testing"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Equals(tt.args.columnName, tt.args.value))
		})
	}
}

func TestNotEquals(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      "testing",
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "!=",
				Values:     []any{"testing"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NotEquals(tt.args.columnName, tt.args.value))
		})
	}
}

func TestGraterThan(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      42,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   ">",
				Values:     []any{42},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GraterThan(tt.args.columnName, tt.args.value))
		})
	}
}

func TestGreaterThanEqual(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      42,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   ">=",
				Values:     []any{42},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GreaterThanEqual(tt.args.columnName, tt.args.value))
		})
	}
}

func TestLessThan(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      42,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "<",
				Values:     []any{42},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, LessThan(tt.args.columnName, tt.args.value))
		})
	}
}

func TestLessThanEqual(t *testing.T) {
	type args struct {
		columnName string
		value      any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      42,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "<=",
				Values:     []any{42},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, LessThanEqual(tt.args.columnName, tt.args.value))
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		columnName string
		value      []any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      []any{42, 56, 127},
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "IN",
				Values:     []any{42, 56, 127},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, In(tt.args.columnName, tt.args.value))
		})
	}
}

func TestNotIn(t *testing.T) {
	type args struct {
		columnName string
		value      []any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				value:      []any{42, 56, 127},
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "NOT IN",
				Values:     []any{42, 56, 127},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NotIn(tt.args.columnName, tt.args.value))
		})
	}
}

func TestBetween(t *testing.T) {
	type args struct {
		columnName string
		val1       any
		val2       any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				val1:       42,
				val2:       56,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "BETWEEN",
				Values:     []any{42, 56},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Between(tt.args.columnName, tt.args.val1, tt.args.val2))
		})
	}
}

func TestNotBetween(t *testing.T) {
	type args struct {
		columnName string
		val1       any
		val2       any
	}
	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
				val1:       42,
				val2:       56,
			},
			want: intypes.SimpleCondition{
				ColumnName: "col1",
				Operator:   "NOT BETWEEN",
				Values:     []any{42, 56},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NotBetween(tt.args.columnName, tt.args.val1, tt.args.val2))
		})
	}
}

func TestGroupedAnd(t *testing.T) {
	type args struct {
		cond1           intypes.Condition
		cond2           intypes.Condition
		additionalConds []intypes.Condition
	}

	testCond1 := LessThanEqual("t1.col1", 42)
	testCond2 := GraterThan("col2", 56)
	testCond3 := Equals("col3", "testing")

	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success; Minimal",
			args: args{
				cond1: testCond1,
				cond2: testCond2,
			},
			want: intypes.GroupedConditions{
				Conjunction: " AND ",
				Conditions:  []intypes.Condition{testCond1, testCond2},
			},
		},
		{
			name: "Success; Additional Conds",
			args: args{
				cond1:           testCond1,
				cond2:           testCond2,
				additionalConds: []intypes.Condition{testCond3},
			},
			want: intypes.GroupedConditions{
				Conjunction: " AND ",
				Conditions:  []intypes.Condition{testCond1, testCond2, testCond3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GroupedAnd(tt.args.cond1, tt.args.cond2, tt.args.additionalConds...))
		})
	}
}

func TestGroupedOr(t *testing.T) {
	type args struct {
		cond1           intypes.Condition
		cond2           intypes.Condition
		additionalConds []intypes.Condition
	}

	testCond1 := LessThanEqual("t1.col1", 42)
	testCond2 := GraterThan("col2", 56)
	testCond3 := Equals("col3", "testing")

	tests := []struct {
		name string
		args args
		want intypes.Condition
	}{
		{
			name: "Success; Minimal",
			args: args{
				cond1: testCond1,
				cond2: testCond2,
			},
			want: intypes.GroupedConditions{
				Conjunction: " OR ",
				Conditions:  []intypes.Condition{testCond1, testCond2},
			},
		},
		{
			name: "Success; Additional Conds",
			args: args{
				cond1:           testCond1,
				cond2:           testCond2,
				additionalConds: []intypes.Condition{testCond3},
			},
			want: intypes.GroupedConditions{
				Conjunction: " OR ",
				Conditions:  []intypes.Condition{testCond1, testCond2, testCond3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GroupedOr(tt.args.cond1, tt.args.cond2, tt.args.additionalConds...))
		})
	}
}

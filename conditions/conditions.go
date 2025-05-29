package conds

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
)

func ColumnValue(columnName string) inconds.ColumnValue {
	panic("unimplemented")
}

func Equals(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "=",
		Values:     []any{value},
	}
}

func NotEquals(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "!=",
		Values:     []any{value},
	}
}

func GraterThan(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">",
		Values:     []any{value},
	}
}

func GreaterThanEqual(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">=",
		Values:     []any{value},
	}
}

func LessThan(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<",
		Values:     []any{value},
	}
}

func LessThanEqual(columnName string, value any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<=",
		Values:     []any{value},
	}
}

func In(columnName string, value []any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "IN",
		Values:     value,
	}
}

func NotIn(columnName string, value []any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT IN",
		Values:     value,
	}
}

func Between(columnName string, val1, val2 any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "BETWEEN",
		Values:     []any{val1, val2},
	}
}

func NotBetween(columnName string, val1, val2 any) inconds.Condition {
	return inconds.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT BETWEEN",
		Values:     []any{val1, val2},
	}
}

func GroupedAnd(cond1, cond2 inconds.Condition, additionalConds ...inconds.Condition) inconds.Condition {
	conds := []inconds.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return inconds.GroupedConditions{
		Conjunction: "AND",
		Conditions:  conds,
	}
}

func GroupedOr(cond1, cond2 inconds.Condition, additionalConds ...inconds.Condition) inconds.Condition {
	conds := []inconds.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return inconds.GroupedConditions{
		Conjunction: "OR",
		Conditions:  conds,
	}
}

package condition

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
)

func ColumnValue(columnName string) incondition.ColumnValue {
	return incondition.ColumnValue{
		ColumnName: columnName,
	}
}

func Equals(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "=",
		Values:     []any{value},
	}
}

func NotEquals(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "!=",
		Values:     []any{value},
	}
}

func GreaterThan(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">",
		Values:     []any{value},
	}
}

func GreaterThanEqual(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">=",
		Values:     []any{value},
	}
}

func LessThan(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<",
		Values:     []any{value},
	}
}

func LessThanEqual(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<=",
		Values:     []any{value},
	}
}

func In(columnName string, value []any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "IN",
		Values:     value,
	}
}

func NotIn(columnName string, value []any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT IN",
		Values:     value,
	}
}

func Between(columnName string, val1, val2 any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "BETWEEN",
		Values:     []any{val1, val2},
	}
}

func NotBetween(columnName string, val1, val2 any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT BETWEEN",
		Values:     []any{val1, val2},
	}
}

func GroupedAnd(cond1, cond2 incondition.Condition, additionalConds ...incondition.Condition) incondition.Condition {
	conds := []incondition.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return incondition.GroupedConditions{
		Conjunction: "AND",
		Conditions:  conds,
	}
}

func GroupedOr(cond1, cond2 incondition.Condition, additionalConds ...incondition.Condition) incondition.Condition {
	conds := []incondition.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return incondition.GroupedConditions{
		Conjunction: "OR",
		Conditions:  conds,
	}
}

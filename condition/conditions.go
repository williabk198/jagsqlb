package condition

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
)

// ColumnValue is to be used in a condition to represent a value from a table column.
//
// For example, if you wanted to create a condition like `t1.col1 > t2.col2` then
// you can use `ColumnValue` like so:
//
//	condition.GreaterThan("t1.col1", condition.ColumnValue("t2.col2"))
//
// This will indicate to not parameterize the value when `Parameterize` is called
// on the condition this ColumnValue is a part of.
func ColumnValue(columnName string) incondition.ColumnValue {
	return incondition.ColumnValue{
		ColumnName: columnName,
	}
}

// Equals returns a condition that can be used in building `WHERE` and `JOIN` clauses that equates a column to a value
func Equals(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "=",
		Values:     []any{value},
	}
}

// NotEquals returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should not be equal to the given value
func NotEquals(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "!=",
		Values:     []any{value},
	}
}

// GreaterThan returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should be greater than the given value
func GreaterThan(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">",
		Values:     []any{value},
	}
}

// GreaterThanEqual returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should be greater than or equal to the given value
func GreaterThanEqual(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   ">=",
		Values:     []any{value},
	}
}

// LessThan returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should be less than the given value
func LessThan(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<",
		Values:     []any{value},
	}
}

// LessThanEqual returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should be less than or equal to the given value
func LessThanEqual(columnName string, value any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "<=",
		Values:     []any{value},
	}
}

// IsNull returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should be NULL
func IsNull(columnName string) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "IS",
		Values:     []any{"NULL"},
	}
}

// IsNotNull returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a column should not be NULL
func IsNotNull(columnName string) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "IS NOT",
		Values:     []any{"NULL"},
	}
}

// In returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a columns value should be in the provided slice of values
func In(columnName string, value []any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "IN",
		Values:     value,
	}
}

// NotIn returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a columns value should not be in the provided slice of values
func NotIn(columnName string, value []any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT IN",
		Values:     value,
	}
}

// Between returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a columns value should be between the two provided values
func Between(columnName string, val1, val2 any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "BETWEEN",
		Values:     []any{val1, val2},
	}
}

// NotBetween returns a condition that can be used in building `WHERE` and `JOIN` clauses that
// indicates a columns value should not be between the two provided values
func NotBetween(columnName string, val1, val2 any) incondition.Condition {
	return incondition.SimpleCondition{
		ColumnName: columnName,
		Operator:   "NOT BETWEEN",
		Values:     []any{val1, val2},
	}
}

// GroupedAnd returns a grouping of conditions with the AND operator between each of them. This can be used like any other condition
// For example if you wanted to create the condition `(t1.col1 = 42 AND t1.col2 != 'testing')`, then you can use `GroupedAnd` like so:
//
//	condition.GroupedAnd(condition.Equals("t1.col1", 42), condition.NotEqual("t1.col2", "testing"))
func GroupedAnd(cond1, cond2 incondition.Condition, additionalConds ...incondition.Condition) incondition.Condition {
	conds := []incondition.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return incondition.GroupedConditions{
		Conjunction: "AND",
		Conditions:  conds,
	}
}

// GroupedOr returns a grouping of conditions with the OR operator between each of them. This can be used like any other condition
// For example if you wanted to create the condition `(t1.col1 = 42 OR t1.col2 < 55)`, then you can use `GroupedOr` like so:
//
//	condition.GroupedOr(condition.Equals("t1.col1", 42), condition.LessThan("t1.col2", 55))
func GroupedOr(cond1, cond2 incondition.Condition, additionalConds ...incondition.Condition) incondition.Condition {
	conds := []incondition.Condition{cond1, cond2}
	conds = append(conds, additionalConds...)

	return incondition.GroupedConditions{
		Conjunction: "OR",
		Conditions:  conds,
	}
}

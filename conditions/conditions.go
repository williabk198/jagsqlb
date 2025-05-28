package conds

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
)

func Equals(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func NotEquals(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func GraterThan(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func GreaterThanEqual(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func LessThan(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func LessThanEqual(columnName string, value any) inconds.Condition {
	panic("unimplemented")
}

func In(columnName string, value []any) inconds.Condition {
	panic("unimplemented")
}

func NotIn(columnName string, value []any) inconds.Condition {
	panic("unimplemented")
}

func Between(columnName string, val1, val2 any) inconds.Condition {
	panic("unimplemented")
}

func NotBetween(columnName string, val1, val2 any) inconds.Condition {
	panic("unimplemented")
}

func GroupedAnd(cond1, cond2 inconds.Condition, additionalConds ...inconds.Condition) inconds.Condition {
	panic("unimplemented")
}

func GroupedOr(cond1, cond2 inconds.Condition, additionalConds ...inconds.Condition) inconds.Condition {
	panic("unimplemented")
}

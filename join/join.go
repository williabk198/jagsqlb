package join

import (
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	injoin "github.com/williabk198/jagsqlb/internal/join"
)

const (
	TypeInner injoin.Type = "INNER JOIN"
	TypeOuter injoin.Type = "OUTER JOIN"
	TypeLeft  injoin.Type = "LEFT JOIN"
	TypeRight injoin.Type = "RIGHT JOIN"
	TypeCross injoin.Type = "CROSS JOIN"
)

// On represent the "ON" portion of a "JOIN" clause and defines how the table will be joined based on the provided conditions.
// If multiple conditions are provided, "AND" will be the operator between the conditions.
func On(condition incondition.Condition, additionalConds ...incondition.Condition) injoin.Relation {
	return injoin.Relation{
		Keyword:  "ON",
		Relation: append([]incondition.Condition{condition}, additionalConds...),
	}
}

// Using represents the "USING" part of a "JOIN" caluse. If the tables being joined both have the provided column name,
// then it will be joined using that column name.
func Using(columnName string) injoin.Relation {
	return injoin.Relation{
		Keyword:  "USING",
		Relation: columnName,
	}
}

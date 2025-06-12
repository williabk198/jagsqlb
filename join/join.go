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

func On(condition incondition.Condition, additionalConds ...incondition.Condition) injoin.Relation {
	return injoin.Relation{
		Keyword:  "ON",
		Relation: append([]incondition.Condition{condition}, additionalConds...),
	}
}

func Using(columnName string) injoin.Relation {
	return injoin.Relation{
		Keyword:  "USING",
		Relation: columnName,
	}
}

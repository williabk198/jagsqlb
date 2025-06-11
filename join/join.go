package join

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
)

const (
	TypeInner injoin.Type = "INNER JOIN"
	TypeOuter injoin.Type = "OUTER JOIN"
	TypeLeft  injoin.Type = "LEFT JOIN"
	TypeRight injoin.Type = "RIGHT JOIN"
	TypeCross injoin.Type = "CROSS JOIN"
)

func On(condition inconds.Condition, additionalConds ...inconds.Condition) injoin.Relation {
	return injoin.Relation{
		Keyword:  "ON",
		Relation: append([]inconds.Condition{condition}, additionalConds...),
	}
}

func Using(columnName string) injoin.Relation {
	return injoin.Relation{
		Keyword:  "USING",
		Relation: columnName,
	}
}

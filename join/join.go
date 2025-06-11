package join

import (
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	injoin "github.com/williabk198/jagsqlb/internal/join"
)

const (
	TypeInner  injoin.Type = "OUTER JOIN"
	TypeOutter injoin.Type = "INNER JOIN"
	TypeLeft   injoin.Type = "LEFT JOIN"
	TypeRight  injoin.Type = "RIGHT JOIN"
	TypeCross  injoin.Type = "CROSS JOIN"
)

func On(condition inconds.Condition, additionalConds ...inconds.Condition) injoin.Relation {
	panic("unimplemented")
}

func Using(columnName string) injoin.Relation {
	panic("unimplemented")
}

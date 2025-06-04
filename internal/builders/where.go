package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	inconds "github.com/williabk198/jagsqlb/internal/conditions"
	"github.com/williabk198/jagsqlb/types"
)

type whereBuilder struct {
	mainQuery  builders.Builder
	conditions []whereCondition
}

func (w whereBuilder) Build() (query string, queryParams []any, err error) {
	sb := new(strings.Builder)
	mainQueryStr, params, err := w.mainQuery.Build()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build preceding query: %w", err)
	}

	sb.WriteString(mainQueryStr[:len(mainQueryStr)-1]) // write the primary query string without the trailing ";"
	sb.WriteString(" WHERE ")

	for _, cond := range w.conditions {
		condStr, condParams, err := cond.condition.Parameterize()
		if err != nil {
			return "", nil, fmt.Errorf("failed to parameterize condition %q: %w", cond, err)
		}

		params = append(params, condParams...)
		if cond.conjunction != "" {
			sb.WriteRune(' ')
			sb.WriteString(cond.conjunction)
			sb.WriteRune(' ')
		}
		sb.WriteString(condStr)
	}
	sb.WriteRune(';')

	return finalizeQuery(sb.String()), params, nil
}

func (w whereBuilder) And(cond inconds.Condition, additionalConds ...inconds.Condition) builders.WhereBuilder {

	w.conditions = append(w.conditions, whereCondition{
		conjunction: "AND",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		w.conditions = append(w.conditions, whereCondition{
			conjunction: "AND",
			condition:   cond,
		})
	}

	return w
}

func (w whereBuilder) Or(cond inconds.Condition, additionalConds ...inconds.Condition) builders.WhereBuilder {
	w.conditions = append(w.conditions, whereCondition{
		conjunction: "OR",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		w.conditions = append(w.conditions, whereCondition{
			conjunction: "OR",
			condition:   cond,
		})
	}

	return w
}

type whereCondition struct {
	conjunction string
	condition   inconds.Condition
}

// Limit implements builders.WhereBuilder.
func (w whereBuilder) Limit(uint) builders.Builder {
	panic("unimplemented")
}

// Offset implements builders.WhereBuilder.
func (w whereBuilder) Offset(uint) builders.OffsetBuilder {
	panic("unimplemented")
}

// OrderBy implements builders.WhereBuilder.
func (w whereBuilder) OrderBy(types.ColumnOrdering, ...types.ColumnOrdering) builders.OrderByBuilder {
	panic("unimplemented")
}

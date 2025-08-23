package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	"github.com/williabk198/jagsqlb/types"
)

// selectWhereBuilder implements `builders.SelectWhereBuilder` and represents the WHERE clause in a SELECT statement
type selectWhereBuilder struct {
	mainQuery      builders.Builder
	conditions     whereConditions
	existingParams int
}

func (w selectWhereBuilder) Build() (query string, queryParams []any, err error) {
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

	return finalizeQuery(sb.String(), w.existingParams), params, nil
}

func (w selectWhereBuilder) And(cond incondition.Condition, additionalConds ...incondition.Condition) builders.SelectWhereBuilder {

	w.conditions.Append(whereCondition{
		conjunction: "AND",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		w.conditions.Append(whereCondition{
			conjunction: "AND",
			condition:   cond,
		})
	}

	return w
}

func (w selectWhereBuilder) Or(cond incondition.Condition, additionalConds ...incondition.Condition) builders.SelectWhereBuilder {
	w.conditions.Append(whereCondition{
		conjunction: "OR",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		w.conditions.Append(whereCondition{
			conjunction: "OR",
			condition:   cond,
		})
	}

	return w
}

// Limit implements builders.WhereBuilder.
func (w selectWhereBuilder) Limit(limit uint) builders.Builder {
	return limitBuilder{
		precedingBuilder: w,
		limit:            limit,
	}
}

// Offset implements builders.WhereBuilder.
func (w selectWhereBuilder) Offset(offset uint) builders.LimitBuilder {
	return offsetBuilder{
		precedingBuilder: w,
		offset:           offset,
	}
}

// OrderBy implements builders.WhereBuilder.
func (w selectWhereBuilder) OrderBy(ordering types.ColumnOrdering, moreOrderings ...types.ColumnOrdering) builders.OffsetBuilder {
	return orderByBuilder{
		precedingBuilder: w,
		columnOrderings:  append([]types.ColumnOrdering{ordering}, moreOrderings...),
	}
}

//TODO: Look into a better way of handling returningWhereBuilder. A lot of duplicated code here

type returningWhereBuilder struct {
	mainQuery      builders.Builder
	conditions     whereConditions
	existingParams int
}

func (rwb returningWhereBuilder) Build() (query string, queryParams []any, err error) {
	sb := new(strings.Builder)
	mainQueryStr, params, err := rwb.mainQuery.Build()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build preceding query: %w", err)
	}

	sb.WriteString(mainQueryStr[:len(mainQueryStr)-1]) // write the primary query string without the trailing ";"
	sb.WriteString(" WHERE ")

	for _, cond := range rwb.conditions {
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

	return finalizeQuery(sb.String(), rwb.existingParams), params, nil
}

func (rwb returningWhereBuilder) And(cond incondition.Condition, additionalConds ...incondition.Condition) builders.ReturningWhereBuilder {
	rwb.conditions.Append(whereCondition{
		conjunction: "AND",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		rwb.conditions.Append(whereCondition{
			conjunction: "AND",
			condition:   cond,
		})
	}

	return rwb
}

func (rwb returningWhereBuilder) Or(cond incondition.Condition, additionalConds ...incondition.Condition) builders.ReturningWhereBuilder {
	rwb.conditions.Append(whereCondition{
		conjunction: "OR",
		condition:   cond,
	})

	for _, cond := range additionalConds {
		rwb.conditions.Append(whereCondition{
			conjunction: "OR",
			condition:   cond,
		})
	}

	return rwb
}

func (rwb returningWhereBuilder) Returning(column string, moreColumns ...string) builders.Builder {
	rb := returningBuilder{
		prevBuilder: rwb,
	}
	return rb.Returning(column, moreColumns...)
}

type whereCondition struct {
	conjunction string
	condition   incondition.Condition
}

type whereConditions []whereCondition

func (wc *whereConditions) Append(condition whereCondition) {
	*wc = append(*wc, condition)
}

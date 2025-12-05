package inbuilders

import (
	"fmt"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type returningBuilder struct {
	prevBuilder      builders.Builder
	returningColumns []intypes.Column
	errs             intypes.ErrorSlice
}

func (rb returningBuilder) Build() (string, []any, error) {
	if len(rb.errs) > 0 {
		return "", nil, rb.errs
	}

	query, params, err := rb.prevBuilder.Build()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build section before the returning builder: %w", err)
	}

	sb := new(strings.Builder)
	sb.WriteString(query[:len(query)-1])

	if len(rb.returningColumns) > 0 {
		sb.WriteString(" RETURNING ")
		sb.WriteString(rb.returningColumns[0].String())
		for i := 1; i < len(rb.returningColumns); i++ {
			sb.WriteString(", ")
			sb.WriteString(rb.returningColumns[i].String())
		}
	}
	sb.WriteRune(';')

	return sb.String(), params, nil
}

func (rb returningBuilder) Returning(column string, moreColumns ...string) builders.Builder {
	col, err := columnParser.Parse(column)
	if err != nil {
		rb.errs = append(rb.errs, err)
		return rb
	}
	rb.returningColumns = append(rb.returningColumns, col)

	for _, mc := range moreColumns {
		col, err = columnParser.Parse(mc)
		if err != nil {
			rb.errs = append(rb.errs, err)
			continue
		}
		rb.returningColumns = append(rb.returningColumns, col)
	}

	return rb
}

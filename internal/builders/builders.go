package inbuilders

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/williabk198/jagsqlb/builders"
	"github.com/williabk198/jagsqlb/internal/utilities/parsers"
	"github.com/williabk198/jagsqlb/types"
)

var (
	tableParser        = parsers.NewTableParser()
	selectColumnParser = parsers.NewSelectColumnParser()
)

type orderByBuilder struct {
	precedingBuilder builders.Builder
	columnOrderings  []types.ColumnOrdering
}

func (obb orderByBuilder) Build() (string, []any, error) {
	query, params, err := obb.precedingBuilder.Build()
	if err != nil {
		return "", nil, err
	}

	sb := new(strings.Builder)
	ordering, err := obb.columnOrderings[0].Stringify()
	if err != nil {
		return "", nil, err
	}
	sb.WriteString(ordering)

	for i := 1; i < len(obb.columnOrderings); i++ {
		sb.WriteString(", ")
		ordering, err := obb.columnOrderings[i].Stringify()
		if err != nil {
			return "", nil, err
		}
		sb.WriteString(ordering)
	}

	query = fmt.Sprintf("%s ORDER BY %s;", query[:len(query)-1], sb.String())
	return query, params, nil
}

func (oob orderByBuilder) Offset(offset uint) builders.OffsetBuilder {
	return offsetBuilder{
		precedingBuilder: oob,
		offset:           offset,
	}
}

func (oob orderByBuilder) Limit(limit uint) builders.Builder {
	return limitBuilder{
		precedingBuilder: oob,
		limit:            limit,
	}
}

type offsetBuilder struct {
	precedingBuilder builders.Builder
	offset           uint
}

func (ob offsetBuilder) Build() (string, []any, error) {
	query, params, err := ob.precedingBuilder.Build()
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("%s OFFSET %d;", query[:len(query)-1], ob.offset)
	return query, params, nil
}

func (ob offsetBuilder) Limit(limit uint) builders.Builder {
	return limitBuilder{
		precedingBuilder: ob,
		limit:            limit,
	}
}

type limitBuilder struct {
	precedingBuilder builders.Builder
	limit            uint
}

func (lb limitBuilder) Build() (string, []any, error) {
	query, params, err := lb.precedingBuilder.Build()
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("%s OFFSET %d;", query[:len(query)-1], lb.limit)
	return query, params, nil
}

// finalizeQuery replaces any "?" characters in the provided query with "$n" characters
func finalizeQuery(query string) string {
	pattern := regexp.MustCompile(`\?`)
	count := 0
	result := pattern.ReplaceAllStringFunc(query, func(value string) string {
		count++
		return fmt.Sprintf("$%d", count)
	})

	return result
}

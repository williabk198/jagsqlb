package inbuilders

import (
	"fmt"
	"regexp"

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
	panic("unimplemented")
}

func (oob orderByBuilder) Offset(offset uint) builders.OffsetBuilder {
	panic("unimplemented")
}

func (oob orderByBuilder) Limit(limit uint) builders.Builder {
	panic("unimplemented")
}

type offsetBuilder struct {
	precedingBuilder builders.Builder
	offset           uint
}

func (ob offsetBuilder) Build() (string, []any, error) {
	panic("unimplemented")
}

func (ob offsetBuilder) Limit(limit uint) builders.Builder {
	panic("unimplemented")
}

type limitBuilder struct {
	precedingBuilder builders.Builder
	limit            uint
}

func (lb limitBuilder) Build() (string, []any, error) {
	panic("unimplemented")
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

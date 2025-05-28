package inbuilders

import (
	"fmt"
	"regexp"

	"github.com/williabk198/jagsqlb/internal/utilities/parsers"
)

var (
	tableParser        = parsers.NewTableParser()
	selectColumnParser = parsers.NewSelectColumnParser()
)

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

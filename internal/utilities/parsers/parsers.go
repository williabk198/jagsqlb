package parsers

import (
	"fmt"
	"regexp"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type Parser[T any] interface {
	Parse(string) (T, error)
}

var (
	extraWhitespaceRegex = regexp.MustCompile(`\s{2,}`)
)

func sanitizeInput(input string) string {
	// strip out all quotes
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.TrimSpace(input)
	// replace all consecutive whitespace characters with a single space character
	return extraWhitespaceRegex.ReplaceAllString(input, " ")
}

func getTableData(input *string) (*intypes.Table, error) {
	var schema string
	splitInput := strings.Split(*input, ".")
	if splitLen := len(splitInput); splitLen == 2 {
		schema = strings.TrimSpace(splitInput[0])
		if schema == "" {
			return nil, intypes.ErrMissingSchemaName
		}

		*input = strings.TrimSpace(splitInput[1])
	} else if splitLen > 2 {
		return nil, intypes.NewInvalidSyntaxError("too many '.' characters")
	}

	if *input == "" {
		return nil, intypes.ErrMissingTableName
	}

	return &intypes.Table{
		Name:   *input,
		Schema: schema,
	}, nil
}

func getAlias(input string) (alias, remainder string, err error) {
	separators := []string{"AS", " "} // "AS" needs to be first here. If " " is first, then a false error will be returned for `multiple occurrences of " "`
	remainder = input

	// attempt to parse out the alias using the pre-defined separators
	for _, sep := range separators {
		splitInput := strings.Split(input, sep)
		if splitLen := len(splitInput); splitLen == 2 {
			remainder = strings.TrimSpace(splitInput[0])
			if remainder == "" {
				return "", "", intypes.NewInvalidSyntaxError("value before alias definition is empty or whitespace")
			}

			alias = strings.TrimSpace(splitInput[1])
			if alias == "" {
				return "", "", intypes.ErrMissingAliasName
			}
		} else if splitLen > 2 {
			return "", "", intypes.NewInvalidSyntaxError(fmt.Sprintf("multiple occurrences of %q", sep))
		}

		if alias != "" {
			break
		}
	}

	if alias != "" && strings.Contains(alias, " ") {
		return "", "", intypes.NewInvalidSyntaxError(fmt.Sprintf("invalid whitespace in alias name: %q", alias))
	}

	if strings.Contains(remainder, " ") {
		return "", "", intypes.NewInvalidSyntaxError(fmt.Sprintf("value before alias definition contains white space: %q", remainder))
	}

	return alias, remainder, nil
}

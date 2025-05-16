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
			return nil, intypes.NewInvalidSytaxError("schema name was not provided")
		}

		*input = strings.TrimSpace(splitInput[1])
	} else if splitLen > 2 {
		return nil, intypes.NewInvalidSytaxError("too many '.' characters")
	}

	if *input == "" {
		return nil, intypes.NewInvalidSytaxError("table name was not provided")
	}

	return &intypes.Table{
		Name:   *input,
		Schema: schema,
	}, nil
}

func getAlias(input *string, seperator string) (string, error) {
	var alias string
	splitInput := strings.Split(*input, seperator)
	if splitLen := len(splitInput); splitLen == 2 {
		*input = strings.TrimSpace(splitInput[0])
		if *input == "" {
			return "", intypes.NewInvalidSytaxError("value before alias definition is empty or whitespace")
		}

		alias = strings.TrimSpace(splitInput[1])
		if alias == "" {
			return "", intypes.NewInvalidSytaxError("alias name not provided")
		}
	} else if splitLen > 2 {
		return "", intypes.NewInvalidSytaxError(fmt.Sprintf("multiple occurences of %q", seperator))
	}

	return alias, nil
}

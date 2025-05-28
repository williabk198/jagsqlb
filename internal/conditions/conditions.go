package inconds

import "github.com/williabk198/jagsqlb/internal/utilities/parsers"

var (
	columnParser = parsers.NewColumnParser()
)

// Condition is essentially the same as a Builder, but needed a way to differentiate a Condtion from a Builder
type Condition interface {
	Parameterize() (string, []any, error)
}

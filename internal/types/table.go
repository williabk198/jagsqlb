package intypes

import (
	"fmt"
	"strings"
)

type Table struct {
	Alias  string
	Name   string
	Schema string
}

// ReferenceString returns a string that can be used as a reference to the table.
// If the table has an alias, then this returns said alias. If the table has a schema defined,
// then the both the schema and the table name will be returned in the format "schema.table".
// Otherwise, just the Name field is returned.
func (t Table) ReferenceString() string {
	refStr := t.Name
	if t.Alias != "" {
		refStr = t.Alias
	}

	if t.Schema == "" || t.Alias != "" {
		return fmt.Sprintf("%q", refStr)
	}

	return fmt.Sprintf("%q.%q", t.Schema, t.Name)

}

func (t Table) String() string {
	sb := new(strings.Builder)

	if t.Schema != "" {
		sb.WriteRune('"')
		sb.WriteString(t.Schema)
		sb.WriteString(`".`)
	}

	sb.WriteRune('"')
	sb.WriteString(t.Name)
	sb.WriteRune('"')

	if t.Alias != "" {
		sb.WriteString(` AS "`)
		sb.WriteString(t.Alias)
		sb.WriteRune('"')
	}

	return sb.String()
}

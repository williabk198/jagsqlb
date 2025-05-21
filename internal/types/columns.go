package intypes

import "strings"

type Column struct {
	Name  string
	Table *Table
}

func (c Column) String() string {
	sb := new(strings.Builder)

	if c.Table != nil {
		// Don't use Table's String method since that can return its alias which is unwanted here.
		sb.WriteString(c.Table.ReferenceString())
		sb.WriteRune('.')
	}

	if c.Name == "*" {
		sb.WriteString(c.Name)
	} else {
		sb.WriteRune('"')
		sb.WriteString(c.Name)
		sb.WriteRune('"')
	}

	return sb.String()
}

type SelectColumn struct {
	Alias string
	Column
}

func (sc SelectColumn) String() string {
	result := sc.Column.String()
	if sc.Alias == "" || sc.Name == "*" {
		return result
	}

	sb := new(strings.Builder)
	sb.WriteString(result)
	sb.WriteString(` AS "`)
	sb.WriteString(sc.Alias)
	sb.WriteRune('"')

	return sb.String()
}

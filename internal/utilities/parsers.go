package inutilities

import intypes "github.com/williabk198/jagsqlb/internal/types"

func ParseTableData(input string) (intypes.Table, error) {
	panic("unimplemented")
}

func ParseColumnData(input string) (intypes.Column, error) {
	panic("unimplemented")
}

func ParseSelectorColumnData(input string) (intypes.SelectorColumn, error) {
	// TODO: Use ParseColumnData and then parse out the column alias if it exsists
	panic("unimplemented")
}

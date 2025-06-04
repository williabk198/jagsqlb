package types

type ordering string

const (
	OrderingAscending  ordering = "ASC"
	OrderingDescending ordering = "DESC"
)

type ColumnOrdering struct {
	ColumnName string
	Ordering   ordering
}

func (co ColumnOrdering) Stringify() (string, error) {
	panic("unimplemented")
}

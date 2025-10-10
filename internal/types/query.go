package intypes

type QueryMarshaler interface {
	MarshalQuery() (string, error)
}

type QueryType string

const (
	QueryTypeInsert QueryType = "insert"
	QueryTypeUpdate QueryType = "update"
)

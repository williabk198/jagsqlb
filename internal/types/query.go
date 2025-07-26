package intypes

type QueryMarshaler interface {
	MarshalQuery() (string, error)
}

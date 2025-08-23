package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnOrdering_Stringify(t *testing.T) {
	tests := []struct {
		name      string
		co        ColumnOrdering
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			co: ColumnOrdering{
				ColumnName: "col1",
				Ordering:   OrderingAscending,
			},
			want:      `"col1" ASC`,
			assertion: assert.NoError,
		},
		{
			name: "Error",
			co: ColumnOrdering{
				ColumnName: ".col1",
				Ordering:   OrderingDescending,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.co.Stringify()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

package intypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_ReferenceString(t *testing.T) {
	tests := []struct {
		name string
		tr   Table
		want string
	}{
		{
			name: "Table With Only Name",
			tr: Table{
				Name: "testTable",
			},
			want: `"testTable"`,
		},
		{
			name: "Table With Schema",
			tr: Table{
				Name:   "testTable",
				Schema: "testing",
			},
			want: `"testing"."testTable"`,
		},
		{
			name: "Table With Alias",
			tr: Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			want: `"tt"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tr.ReferenceString())
		})
	}
}

func TestTable_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Table
		want string
	}{
		{
			name: "Table Only",
			tr: Table{
				Name: "testTable",
			},
			want: `"testTable"`,
		},
		{
			name: "Table with Schema",
			tr: Table{
				Name:   "testTable",
				Schema: "testing",
			},
			want: `"testing"."testTable"`,
		},
		{
			name: "Table with Alias",
			tr: Table{
				Alias: "tt",
				Name:  "testTable",
			},
			want: `"testTable" AS "tt"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tr.String())
		})
	}
}

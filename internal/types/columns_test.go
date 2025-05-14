package intypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumn_String(t *testing.T) {
	tests := []struct {
		name string
		c    Column
		want string
	}{
		{
			name: "Column Only",
			c: Column{
				Name: "testCol",
			},
			want: `"testCol"`,
		},
		{
			name: "Column with Table",
			c: Column{
				Name:  "testCol",
				Table: &Table{Name: "testTable"},
			},
			want: `"testTable"."testCol"`,
		},
		{
			name: "Column with Table and Schema",
			c: Column{
				Name: "testCol",
				Table: &Table{
					Name:   "testTable",
					Schema: "testing",
				},
			},
			want: `"testing"."testTable"."testCol"`,
		},
		{
			name: "Column with Aliased Table",
			c: Column{
				Name: "testCol",
				Table: &Table{
					Alias:  "t",
					Name:   "testTable",
					Schema: "testing",
				},
			},
			want: `"t"."testCol"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}

func TestSelectorColumn_String(t *testing.T) {
	tests := []struct {
		name string
		sc   SelectorColumn
		want string
	}{
		{
			name: "Column Only",
			sc: SelectorColumn{
				Column: Column{
					Name: "testCol",
				},
			},
			want: `"testCol1"`,
		},
		{
			name: "Column with Alias",
			sc: SelectorColumn{
				Alias:  "tc",
				Column: Column{Name: "testCol"},
			},
			want: `"testCol" AS "tc"`,
		},
		{
			name: "Column with Table",
			sc: SelectorColumn{
				Column: Column{
					Name:  "testCol",
					Table: &Table{Name: "testTable"},
				},
			},
			want: `"testTable"."testCol"`,
		},
		{
			name: "Column with Table and Schema",
			sc: SelectorColumn{
				Column: Column{
					Name: "testCol",
					Table: &Table{
						Name:   "testTable",
						Schema: "testing",
					},
				},
			},
			want: `"testing"."testTable"."testCol"`,
		},
		{
			name: "Column with Aliased Table",
			sc: SelectorColumn{
				Column: Column{
					Name: "testCol",
					Table: &Table{
						Alias:  "t",
						Name:   "testTable",
						Schema: "testing",
					},
				},
			},
			want: `"t"."testCol"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sc.String())
		})
	}
}

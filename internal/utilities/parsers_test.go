package inutilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func TestParseTableData(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      intypes.Table
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success; Table Name Only",
			args:      args{input: "testTable"},
			want:      intypes.Table{Name: "testTable"},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema",
			args: args{input: "testing.testTable"},
			want: intypes.Table{
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema with AS Alias",
			args: args{input: "testing.testTable AS tt"},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema with Space Alias",
			args: args{input: "testing.testTable tt"},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Quoted Name, Schema and Alias",
			args: args{input: `"testing"."testTable" AS "tt"`},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Missing Schema",
			args:      args{input: ".testTable"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Name",
			args:      args{input: "testing."},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Alias",
			args:      args{input: "testTable AS"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTableData(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseColumnData(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      intypes.Column
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success; Name Only",
			args:      args{input: "testCol"},
			want:      intypes.Column{Name: "testCol"},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Table",
			args: args{input: "testTable.testCol"},
			want: intypes.Column{
				Name:  "testCol",
				Table: intypes.Table{Name: "testTable"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name, Table and Schema",
			args: args{input: "testing.testTable.testCol"},
			want: intypes.Column{
				Name: "testCol",
				Table: intypes.Table{
					Name:   "testTable",
					Schema: "testing",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Quoted Name, Table and Schema",
			args: args{input: `"testing"."testTable"."testCol"`},
			want: intypes.Column{
				Name: "testCol",
				Table: intypes.Table{
					Name:   "testTable",
					Schema: "testing",
				},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Missing Table Name",
			args:      args{input: " .testCol"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Column Name",
			args:      args{input: "testTable. "},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Table Schema",
			args:      args{input: " .tableName.testCol"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Table with Schema",
			args:      args{input: "testing. .testCol"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColumnData(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseSelectorColumnData(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name      string
		args      args
		want      intypes.SelectorColumn
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; AS Alias",
			args: args{input: "testCol AS tc"},
			want: intypes.SelectorColumn{
				Alias: "tc",
				Column: intypes.Column{
					Name: "testCol",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Space Alias",
			args: args{input: "testing.testTable.testCol tc"},
			want: intypes.SelectorColumn{
				Alias: "tc",
				Column: intypes.Column{
					Name: "testCol",
					Table: intypes.Table{
						Name:   "testTable",
						Schema: "testing",
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Missing Alias",
			args:      args{input: "testCol AS "},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSelectorColumnData(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

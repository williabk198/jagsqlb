package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_tableParser_Parse(t *testing.T) {
	type args struct {
		tableStr string
	}
	tests := []struct {
		name      string
		tp        tableParser
		args      args
		want      intypes.Table
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success; Table Name Only",
			args:      args{tableStr: "testTable"},
			want:      intypes.Table{Name: "testTable"},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema",
			args: args{tableStr: "testing.testTable"},
			want: intypes.Table{
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema with AS Alias",
			args: args{tableStr: "testing.testTable AS tt"},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Schema with Space Alias",
			args: args{tableStr: "testing.testTable tt"},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Quoted Name, Schema and Alias",
			args: args{tableStr: `"testing"."testTable" AS "tt"`},
			want: intypes.Table{
				Alias:  "tt",
				Name:   "testTable",
				Schema: "testing",
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error, Empty Input",
			args:      args{tableStr: ""},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Schema",
			args:      args{tableStr: ".testTable"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Name",
			args:      args{tableStr: "testing."},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Alias",
			args:      args{tableStr: "testTable AS"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tp.Parse(tt.args.tableStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

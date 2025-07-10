package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnTagParser(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name      string
		args      args
		want      map[string]any
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; No Tag",
			args: args{
				input: struct{ Data string }{"testing"},
			},
			want:      map[string]any{"Data": "testing"},
			assertion: assert.NoError,
		},
		{
			name: "Success; With Tag",
			args: args{
				input: struct {
					Data int `jagsqlb:"data"`
				}{52},
			},
			want: map[string]any{
				"data": 52,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Incorrect Parameter Type",
			args: args{
				input: "badInput",
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ColumnTagParser(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

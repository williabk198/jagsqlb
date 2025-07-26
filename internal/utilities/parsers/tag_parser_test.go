package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnTagParser(t *testing.T) {
	type args struct {
		input any
	}
	type wants struct {
		cols []string
		vals []any
	}

	type innerTestStruct struct {
		Data string `jagsqlb:"inner_data"`
	}

	tests := []struct {
		name      string
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; No Tag",
			args: args{
				input: struct{ Data string }{"testing"},
			},
			wants: wants{
				cols: []string{"Data"},
				vals: []any{"testing"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; With Tag",
			args: args{
				input: struct {
					Data int `jagsqlb:"data"`
				}{52},
			},
			wants: wants{
				cols: []string{"data"},
				vals: []any{52},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; With time.Time Value",
			args: args{
				input: struct {
					Timestamp time.Time `jagsqlb:"ts"`
				}{time.Unix(1753502400, 0)},
			},
			wants: wants{
				cols: []string{"ts"},
				vals: []any{time.Unix(1753502400, 0)},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; With Inline Struct",
			args: args{
				input: struct {
					Data      string          `jagsqlb:"data"`
					InnerData innerTestStruct `jagsqlb:";inline"`
				}{
					Data: "outer",
					InnerData: innerTestStruct{
						Data: "inner",
					},
				},
			},
			wants: wants{
				cols: []string{"data", "inner_data"},
				vals: []any{"outer", "inner"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; With Omitted Field",
			args: args{
				input: struct {
					Data      string          `jagsqlb:"data"`
					InnerData innerTestStruct `jagsqlb:";omit"`
				}{
					Data: "outer",
					InnerData: innerTestStruct{
						Data: "inner",
					},
				},
			},
			wants: wants{
				cols: []string{"data"},
				vals: []any{"outer"},
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
			gotCols, gotVals, err := ParseColumnTag(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.cols, gotCols)
			assert.Equal(t, tt.wants.vals, gotVals)
		})
	}
}

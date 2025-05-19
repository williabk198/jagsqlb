package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAlias(t *testing.T) {
	type args struct {
		input string
	}
	type wants struct {
		alias     string
		remainder string
	}
	tests := []struct {
		name      string
		args      args
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success with no Alias",
			args:      args{input: "someVal"},
			wants:     wants{remainder: "someVal"},
			assertion: assert.NoError,
		},
		{
			name: "Success with 'AS'",
			args: args{input: "someVal AS sv"},
			wants: wants{
				alias:     "sv",
				remainder: "someVal",
			},
			assertion: assert.NoError,
		},
		{
			name: "Success with Space",
			args: args{input: "someVal sv"},
			wants: wants{
				alias: "sv", remainder: "someVal",
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; No Data before Alias",
			args:      args{input: " AS sv"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Too Many 'AS' Definitions",
			args:      args{input: "somVal AS sv AS vs"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Too Many Space Definitions",
			args:      args{input: "somVal sv vs"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Partial Definition",
			args:      args{input: "someVal AS"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Space and 'AS' Alias Provided",
			args:      args{input: "someVal c AS t"},
			assertion: assert.Error,
		},
		{
			name:      "Error; 'AS' and Space Alias Provided",
			args:      args{input: "someVal AS c t"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAlias, gotRemainder, err := getAlias(tt.args.input)
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.alias, gotAlias)
			assert.Equal(t, tt.wants.remainder, gotRemainder)
		})
	}
}

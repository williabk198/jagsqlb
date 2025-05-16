package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_columnParser_Parse(t *testing.T) {
	type args struct {
		columnStr string
	}

	testTableParser := tableParser{}
	testColumnParser := columnParser{tableParser: testTableParser}

	tests := []struct {
		name      string
		cp        columnParser
		args      args
		want      intypes.Column
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success; Name Only",
			cp:        testColumnParser,
			args:      args{columnStr: "testCol"},
			want:      intypes.Column{Name: "testCol"},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name and Table",
			cp:   testColumnParser,
			args: args{columnStr: "testTable.testCol"},
			want: intypes.Column{
				Name:  "testCol",
				Table: &intypes.Table{Name: "testTable"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Name, Table and Schema",
			cp:   testColumnParser,
			args: args{columnStr: "testing.testTable.testCol"},
			want: intypes.Column{
				Name: "testCol",
				Table: &intypes.Table{
					Name:   "testTable",
					Schema: "testing",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Quoted Name, Table and Schema",
			cp:   testColumnParser,
			args: args{columnStr: `"testing"."testTable"."testCol"`},
			want: intypes.Column{
				Name: "testCol",
				Table: &intypes.Table{
					Name:   "testTable",
					Schema: "testing",
				},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Empty Input",
			cp:        testColumnParser,
			args:      args{columnStr: ""},
			assertion: assert.Error,
		},
		{
			name:      "Error; Excessive Periods",
			cp:        testColumnParser,
			args:      args{columnStr: "invalid.schema.table.column"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Table Name",
			cp:        testColumnParser,
			args:      args{columnStr: " .testCol"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Column Name",
			cp:        testColumnParser,
			args:      args{columnStr: "testTable. "},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Table Schema",
			cp:        testColumnParser,
			args:      args{columnStr: " .tableName.testCol"},
			assertion: assert.Error,
		},
		{
			name:      "Error; Missing Table with Schema",
			cp:        testColumnParser,
			args:      args{columnStr: "testing. .testCol"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cp.Parse(tt.args.columnStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_selectColumnParser_Parse(t *testing.T) {
	type args struct {
		selectColumnStr string
	}

	testTableParser := tableParser{}
	testSelectColumnParser := selectColumnParser{tableParser: testTableParser}

	tests := []struct {
		name      string
		scp       selectColumnParser
		args      args
		want      intypes.SelectColumn
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; AS Alias",
			scp:  testSelectColumnParser,
			args: args{selectColumnStr: "testCol AS tc"},
			want: intypes.SelectColumn{
				Alias: "tc",
				Column: intypes.Column{
					Name: "testCol",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Space Alias",
			scp:  testSelectColumnParser,
			args: args{selectColumnStr: "testing.testTable.testCol tc"},
			want: intypes.SelectColumn{
				Alias: "tc",
				Column: intypes.Column{
					Name: "testCol",
					Table: &intypes.Table{
						Name:   "testTable",
						Schema: "testing",
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name:      "Error; Missing Alias",
			scp:       testSelectColumnParser,
			args:      args{selectColumnStr: "testCol AS "},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scp.Parse(tt.args.selectColumnStr)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

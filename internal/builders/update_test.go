package inbuilders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	"github.com/williabk198/jagsqlb/condition"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_updateBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}
	tests := []struct {
		name      string
		u         updateBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; w/o From",
			u: updateBuilder{
				table: intypes.Table{Name: "table1"},
				columns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
				},
				vals: []any{"testing", 137},
			},
			wants: wants{
				query:  `UPDATE "table1" SET "col1"=$1, "col2"=$2;`,
				params: []any{"testing", 137},
			},
			assertion: assert.NoError,
		},
		{
			name: "Sucess; w/ From",
			u: updateBuilder{
				table: intypes.Table{Name: "table1"},
				columns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
				},
				vals: []any{
					"testing",
					137,
				},
				fromTables: []intypes.Table{
					{Name: "table2"},
					{Name: "table3"},
				},
			},
			wants: wants{
				query:  `UPDATE "table1" SET "col1"=$1, "col2"=$2 FROM "table2", "table3";`,
				params: []any{"testing", 137},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; ErrorSlice not Empty",
			u: updateBuilder{
				errs: intypes.ErrorSlice{assert.AnError},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotQueryParams, err := tt.u.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotQueryParams)
		})
	}
}

func Test_updateBuilder_SetMap(t *testing.T) {
	type args struct {
		colValMap map[string]any
	}
	tests := []struct {
		name string
		u    updateBuilder
		args args
		want builders.UpdateFromWhereBuilder
	}{
		{
			name: "Success",
			u:    updateBuilder{},
			args: args{
				colValMap: map[string]any{"col1": "something"},
			},
			want: updateBuilder{
				columns: []intypes.Column{{Name: "col1"}},
				vals:    []any{"something"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.SetMap(tt.args.colValMap))
		})
	}
}

func Test_updateBuilder_SetStruct(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name string
		u    updateBuilder
		args args
		want builders.UpdateFromWhereBuilder
	}{
		{
			name: "Success; No Struct Tag",
			u:    updateBuilder{},
			args: args{
				value: struct{ Data string }{"testing"},
			},
			want: updateBuilder{
				columns: []intypes.Column{{Name: "Data"}},
				vals:    []any{"testing"},
			},
		},
		{
			name: "Success; With Struct Tag",
			u:    updateBuilder{},
			args: args{
				value: struct {
					Data int `jagsqlb:"data"`
				}{153},
			},
			want: updateBuilder{
				columns: []intypes.Column{{Name: "data"}},
				vals:    []any{153},
			},
		},
		{
			name: "Error; Invalid Input Type",
			u:    updateBuilder{},
			args: args{
				value: "bad_val",
			},
			want: updateBuilder{
				errs: intypes.ErrorSlice{
					fmt.Errorf(
						"failed to process argument of SetStruct: %w",
						fmt.Errorf("recieved value is not a struct type"),
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.SetStruct(tt.args.value))
		})
	}
}

func Test_updateBuilder_From(t *testing.T) {
	type args struct {
		table      string
		moreTables []string
	}
	tests := []struct {
		name string
		u    updateBuilder
		args args
		want builders.ReturningWhereBuilder
	}{
		{
			name: "Success; Single Table",
			u:    updateBuilder{},
			args: args{
				table: "table2",
			},
			want: returningWhereBuilder{
				mainQuery: updateBuilder{
					fromTables: []intypes.Table{
						{Name: "table2"},
					},
				},
			},
		},
		{
			name: "Success; Multiple Tables",
			u:    updateBuilder{},
			args: args{
				table:      "table2",
				moreTables: []string{"table3"},
			},
			want: returningWhereBuilder{
				mainQuery: updateBuilder{
					fromTables: []intypes.Table{
						{Name: "table2"},
						{Name: "table3"},
					},
				},
			},
		},
		{
			name: "Error; Bad Table Name",
			u:    updateBuilder{},
			args: args{
				table: ".bad_table",
			},
			want: returningWhereBuilder{
				mainQuery: updateBuilder{
					errs: intypes.ErrorSlice{
						fmt.Errorf("failed to parse table data from %q: %w", ".bad_table", intypes.ErrMissingSchemaName),
					},
				},
			},
		},
		{
			name: "Error; Bad MoreTable Name",
			u:    updateBuilder{},
			args: args{
				table:      "table2",
				moreTables: []string{".bad_table"},
			},
			want: returningWhereBuilder{
				mainQuery: updateBuilder{
					fromTables: []intypes.Table{{Name: "table2"}},
					errs: intypes.ErrorSlice{
						fmt.Errorf("failed to parse table data from %q: %w", ".bad_table", intypes.ErrMissingSchemaName),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.From(tt.args.table, tt.args.moreTables...))
		})
	}
}

func Test_updateBuilder_Where(t *testing.T) {
	type args struct {
		cond      incondition.Condition
		moreConds []incondition.Condition
	}
	tests := []struct {
		name string
		u    updateBuilder
		args args
		want builders.ReturningWhereBuilder
	}{
		{
			name: "Success",
			u: updateBuilder{
				table:   intypes.Table{Name: "table1"},
				columns: []intypes.Column{{Name: "col1"}},
				vals:    []any{"test"},
			},
			args: args{
				cond: condition.Equals("col2", 56),
			},
			want: returningWhereBuilder{
				mainQuery: updateBuilder{
					table:   intypes.Table{Name: "table1"},
					columns: []intypes.Column{{Name: "col1"}},
					vals:    []any{"test"},
				},
				conditions: whereConditions{
					{
						conjunction: "AND",
						condition:   condition.Equals("col2", 56),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.u.Where(tt.args.cond, tt.args.moreConds...))
		})
	}
}

func TestNewUpdateBuilder(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name string
		args args
		want builders.UpdateBuilder
	}{
		{
			name: "Success",
			args: args{
				table: "table1",
			},
			want: updateBuilder{
				table: intypes.Table{Name: "table1"},
			},
		},
		{
			name: "Error; Bad Table Name",
			args: args{
				table: ".bad_table",
			},
			want: updateBuilder{
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse table data from %q: %w", ".bad_table", intypes.ErrMissingSchemaName),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewUpdateBuilder(tt.args.table))
		})
	}
}

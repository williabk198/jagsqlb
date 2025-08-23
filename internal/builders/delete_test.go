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

func Test_deleteBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		d         deleteBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; No Using Tables",
			d: deleteBuilder{
				table: "table1",
			},
			wants: wants{
				query: `DELETE FROM "table1";`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; With Using Tables",
			d: deleteBuilder{
				table: "table1",
				usingTables: []intypes.Table{
					{Name: "table2"},
					{Name: "table3"},
				},
			},
			wants: wants{
				query: `DELETE FROM "table1" USING "table2", "table3";`,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; Bad Table Name",
			d: deleteBuilder{
				table: ".table1",
			},
			assertion: assert.Error,
		},
		{
			name: "Error; ErrorSlice not empty",
			d: deleteBuilder{
				table: "table1",
				errs: intypes.ErrorSlice{
					assert.AnError,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotQueryParams, err := tt.d.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, gotQuery)
			assert.Equal(t, tt.wants.params, gotQueryParams)
		})
	}
}

func Test_deleteBuilder_Using(t *testing.T) {
	type args struct {
		tableName string
	}
	tests := []struct {
		name string
		d    deleteBuilder
		args args
		want builders.DeleteBuilder
	}{
		{
			name: "Success",
			d:    deleteBuilder{},
			args: args{
				tableName: "table1",
			},
			want: deleteBuilder{
				usingTables: []intypes.Table{
					{Name: "table1"},
				},
			},
		},
		{
			name: "Error",
			d:    deleteBuilder{},
			args: args{
				tableName: ".table2",
			},
			want: deleteBuilder{
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse table data from %q: %w", ".table2", intypes.ErrMissingSchemaName),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Using(tt.args.tableName))
		})
	}
}

func Test_deleteBuilder_Where(t *testing.T) {
	type args struct {
		condition      incondition.Condition
		moreConditions []incondition.Condition
	}
	tests := []struct {
		name string
		d    deleteBuilder
		args args
		want builders.ReturningWhereBuilder
	}{
		{
			name: "Success",
			d:    deleteBuilder{},
			args: args{
				condition: condition.Equals("col1", 52),
			},
			want: returningWhereBuilder{
				mainQuery: deleteBuilder{},
				conditions: whereConditions{
					{condition: condition.Equals("col1", 52)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Where(tt.args.condition, tt.args.moreConditions...))
		})
	}
}

func Test_deleteBuilder_Returning(t *testing.T) {
	type args struct {
		column      string
		moreColumns []string
	}
	tests := []struct {
		name string
		d    deleteBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success",
			d:    deleteBuilder{},
			args: args{
				column:      "col1",
				moreColumns: []string{"col2"},
			},
			want: returningBuilder{
				prevBuilder: deleteBuilder{},
				returningColumns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Returning(tt.args.column, tt.args.moreColumns...))
		})
	}
}

func TestNewDeleteBuilder(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name string
		args args
		want builders.DeleteBuilder
	}{
		{
			name: "Success",
			args: args{
				table: "table1",
			},
			want: deleteBuilder{
				table: "table1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewDeleteBuilder(tt.args.table))
		})
	}
}

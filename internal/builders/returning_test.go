package inbuilders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/builders"
	"github.com/williabk198/jagsqlb/condition"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func Test_returningBuilder_Build(t *testing.T) {
	type wants struct {
		query  string
		params []any
	}

	tests := []struct {
		name      string
		rb        returningBuilder
		wants     wants
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success; Single Column",
			rb: returningBuilder{
				prevBuilder: NewDeleteBuilder("table1").Where(condition.Equals("col1", "val")),
				returningColumns: []intypes.Column{
					{Name: "*"},
				},
			},
			wants: wants{
				query:  `DELETE FROM "table1" WHERE "col1" = $1 RETURNING *;`,
				params: []any{"val"},
			},
			assertion: assert.NoError,
		},
		{
			name: "Success; Multiple Column",
			rb: returningBuilder{
				prevBuilder: NewDeleteBuilder("table1").Where(condition.GreaterThan("col2", 52)),
				returningColumns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
				},
			},
			wants: wants{
				query:  `DELETE FROM "table1" WHERE "col2" > $1 RETURNING "col1", "col2";`,
				params: []any{52},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error; ErrorSlice not Empty",
			rb: returningBuilder{
				errs: intypes.ErrorSlice{assert.AnError},
			},
			assertion: assert.Error,
		},
		{
			name: "Error; Previous Build Error",
			rb: returningBuilder{
				prevBuilder: NewDeleteBuilder(".table2"),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.rb.Build()
			tt.assertion(t, err)
			assert.Equal(t, tt.wants.query, got)
			assert.Equal(t, tt.wants.params, got1)
		})
	}
}

func Test_returningBuilder_Returning(t *testing.T) {
	type args struct {
		column      string
		moreColumns []string
	}
	tests := []struct {
		name string
		rb   returningBuilder
		args args
		want builders.Builder
	}{
		{
			name: "Success; Single Column",
			rb:   returningBuilder{},
			args: args{
				column: "col1",
			},
			want: returningBuilder{
				returningColumns: []intypes.Column{
					{Name: "col1"},
				},
			},
		},
		{
			name: "Success; Multiple Coulmns",
			rb:   returningBuilder{},
			args: args{
				column:      "col1",
				moreColumns: []string{"col2", "col3"},
			},
			want: returningBuilder{
				returningColumns: []intypes.Column{
					{Name: "col1"},
					{Name: "col2"},
					{Name: "col3"},
				},
			},
		},
		{
			name: "Error; Bad Column Name",
			rb:   returningBuilder{},
			args: args{
				column: ".bad_col",
			},
			want: returningBuilder{
				errs: intypes.ErrorSlice{
					fmt.Errorf("failed to parse table data provided in %q: %w", ".bad_col", intypes.ErrMissingTableName),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rb.Returning(tt.args.column, tt.args.moreColumns...))
		})
	}
}

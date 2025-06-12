package inutilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	intypes "github.com/williabk198/jagsqlb/internal/types"
)

func TestCoalesceSelectColumnsFullString(t *testing.T) {
	type args struct {
		cols []intypes.SelectColumn
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				cols: []intypes.SelectColumn{
					{
						Alias: "t1c1",
						Column: intypes.Column{
							Name: "column1",
							Table: &intypes.Table{
								Name: "table1",
							},
						},
					},
					{
						Column: intypes.Column{
							Name: "column1",
							Table: &intypes.Table{
								Name:   "table1",
								Schema: "metadata",
							},
						},
					},
					{
						Column: intypes.Column{
							Name: "column2",
							Table: &intypes.Table{
								Alias:  "mt2",
								Name:   "table2",
								Schema: "metadata",
							},
						},
					},
					{
						Column: intypes.Column{
							Name:  "*",
							Table: &intypes.Table{Name: "table2"},
						},
					},
				},
			},
			want: `"table1"."column1" AS "t1c1", "metadata"."table1"."column1", "mt2"."column2", "table2".* `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CoalesceSelectColumnsFullString(tt.args.cols))
		})
	}
}

func TestCoalesceSelectColumnNamesString(t *testing.T) {
	type args struct {
		cols []intypes.SelectColumn
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				cols: []intypes.SelectColumn{
					{
						Alias: "t1c1",
						Column: intypes.Column{
							Name: "column1",
							Table: &intypes.Table{
								Name:   "table1",
								Schema: "metadata",
							},
						},
					},
					{
						Column: intypes.Column{
							Name: "column2",
							Table: &intypes.Table{
								Name:   "table1",
								Schema: "metadata",
							},
						},
					},
					{
						Column: intypes.Column{
							Name: "column3",
							Table: &intypes.Table{
								Alias:  "mt2",
								Name:   "table1",
								Schema: "metadata",
							},
						},
					},
				},
			},
			want: `"column1" AS "t1c1", "column2", "column3" `,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CoalesceSelectColumnNamesString(tt.args.cols))
		})
	}
}

func TestCoalesceTablesString(t *testing.T) {
	type args struct {
		tables []intypes.Table
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				tables: []intypes.Table{
					{
						Name:   "table1",
						Schema: "schema",
					},
					{
						Alias: "t2",
						Name:  "table2",
					},
				},
			},
			want: `"schema"."table1", "table2" AS "t2"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CoalesceTablesString(tt.args.tables))
		})
	}
}

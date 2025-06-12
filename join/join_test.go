package join

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williabk198/jagsqlb/condition"
	incondition "github.com/williabk198/jagsqlb/internal/condition"
	injoin "github.com/williabk198/jagsqlb/internal/join"
)

func TestOn(t *testing.T) {
	type args struct {
		condition       incondition.Condition
		additionalConds []incondition.Condition
	}

	testJoinCond1 := condition.Equals("t1.col1", condition.ColumnValue("t2.col2"))
	testJoinCond2 := condition.GreaterThan("t1.col2", 56)

	tests := []struct {
		name string
		args args
		want injoin.Relation
	}{
		{
			name: "Success; Single Condition",
			args: args{
				condition: testJoinCond1,
			},
			want: injoin.Relation{
				Keyword:  "ON",
				Relation: []incondition.Condition{testJoinCond1},
			},
		},
		{
			name: "Success; Multiple Conditions",
			args: args{
				condition:       testJoinCond1,
				additionalConds: []incondition.Condition{testJoinCond2},
			},
			want: injoin.Relation{
				Keyword:  "ON",
				Relation: []incondition.Condition{testJoinCond1, testJoinCond2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, On(tt.args.condition, tt.args.additionalConds...))
		})
	}
}

func TestUsing(t *testing.T) {
	type args struct {
		columnName string
	}
	tests := []struct {
		name string
		args args
		want injoin.Relation
	}{
		{
			name: "Success",
			args: args{
				columnName: "col1",
			},
			want: injoin.Relation{
				Keyword:  "USING",
				Relation: "col1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Using(tt.args.columnName))
		})
	}
}

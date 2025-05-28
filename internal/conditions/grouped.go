package inconds

import (
	"fmt"
	"strings"

	intypes "github.com/williabk198/jagsqlb/internal/types"
)

type GroupedConditions struct {
	Conjunction string // This value should always be either " AND " or " OR "
	Conditions  []Condition
}

func (gc GroupedConditions) Parameterize() (string, []any, error) {
	sb := new(strings.Builder)
	resultParams := make([]any, 0)
	errs := make(intypes.ErrorSlice, 0)

	sb.WriteRune('(')

	str, err := gc.parameterizeHelper(gc.Conditions[0], &resultParams)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parameterize sub-condition %q: %w", gc.Conditions[0], err))
	}

	sb.WriteString(str)

	for i := 1; i < len(gc.Conditions); i++ {
		sb.WriteRune(' ')
		sb.WriteString(gc.Conjunction)
		sb.WriteRune(' ')

		str, err = gc.parameterizeHelper(gc.Conditions[i], &resultParams)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parameterize sub-condition %q: %w", gc.Conditions[i], err))
			continue
		}
		sb.WriteString(str)
	}

	sb.WriteRune(')')

	if len(errs) > 0 {
		return "", nil, errs
	}

	return sb.String(), resultParams, nil
}

func (gc GroupedConditions) parameterizeHelper(cond Condition, currParams *[]any) (string, error) {
	str, params, err := cond.Parameterize()
	if err != nil {
		return "", err
	}

	// If the condition is a grouped condition or the paramaterized string doesn't represent an "IN" condition,
	// then appened each element in `params` to the current slice of query parameters
	if _, ok := cond.(GroupedConditions); ok || !strings.Contains(str, " IN ") {
		*currParams = append(*currParams, params...)
	} else {
		// Otherwise, this condition represents an "IN" SimpleCondition, and we want to append the slice as a single
		// elemnet to the slice of current query parameters
		*currParams = append(*currParams, params)
	}

	return str, nil
}

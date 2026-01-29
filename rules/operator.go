package rules

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrNoMatch = fmt.Errorf("no match")

var operators = map[Operator][]EvalFunc{

	"eq": {
		numeric(func(a, b float64) (bool, error) { return a == b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a == b, nil }),
	},
	"neq": {
		numeric(func(a, b float64) (bool, error) { return a != b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a != b, nil }),
	},
	"gt": {
		numeric(func(a, b float64) (bool, error) { return a > b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a > b, nil }),
	},
	"gte": {
		numeric(func(a, b float64) (bool, error) { return a >= b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a >= b, nil }),
	},
	"lt": {
		numeric(func(a, b float64) (bool, error) { return a < b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a < b, nil }),
	},
	"lte": {
		numeric(func(a, b float64) (bool, error) { return a <= b, nil }),
		stringsOnly(func(a, b string) (bool, error) { return a <= b, nil }),
	},
	"contains": {
		stringsOnly(func(s, sub string) (bool, error) {
			return strings.Contains(s, sub), nil
		}),
		stringSliceString(func(actual string, expected []string) (bool, error) {
			for _, v := range expected {
				if strings.Contains(actual, strings.ToLower(v)) {
					return true, nil
				}
			}
			return false, nil
		}),
	},
	"like": {
		stringsOnly(func(s, re string) (bool, error) {
			return regexp.MatchString(re, s)
		}),
	},
}

type EvalFunc func(actual, expected any) (bool, error)

func numeric(fn func(a, b float64) (bool, error)) EvalFunc {
	return func(a, b any) (bool, error) {
		va, oka := anyToFloat64(a)
		vb, okb := anyToFloat64(b)
		if !oka || !okb {
			return false, ErrNoMatch
		}
		return fn(va, vb)
	}
}

func stringsOnly(fn func(a, b string) (bool, error)) EvalFunc {
	return func(a, b any) (bool, error) {
		sa, ok1 := a.(string)
		sb, ok2 := b.(string)
		if !ok1 || !ok2 {
			return false, ErrNoMatch
		}
		return fn(strings.ToLower(sa), strings.ToLower(sb))
	}
}

func stringSliceString(fn func(a string, b []string) (bool, error)) EvalFunc {
	return func(a, b any) (bool, error) {
		sa, ok := a.(string)
		if !ok {
			return false, ErrNoMatch
		}
		var sliceB []string
		switch b := b.(type) {
		case []any:
			sliceB = make([]string, len(b))
			for i, v := range b {
				vs, ok := v.(string)
				if !ok {
					return false, ErrNoMatch
				}
				sliceB[i] = strings.ToLower(vs)
			}
		case string:
			sliceB = []string{strings.ToLower(b)}
		default:
			return false, ErrNoMatch
		}
		return fn(strings.ToLower(sa), sliceB)
	}
}

func anyToFloat64(v any) (float64, bool) {
	switch x := v.(type) {
	case int:
		return float64(x), true
	case int32:
		return float64(x), true
	case uint32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint64:
		return float64(x), true
	case float32:
		return float64(x), true
	case float64:
		return x, true
	default:
		return 0, false
	}
}

type Operator string

func (o Operator) Evaluate(actual, expected any) (bool, error) {
	evaluators, ok := operators[o]
	if !ok {
		return false, fmt.Errorf("unknown operator %q", o)
	}

	for _, evaluator := range evaluators {
		res, err := evaluator(actual, expected)
		if errors.Is(err, ErrNoMatch) {
			continue
		}
		return res, err
	}

	return false, fmt.Errorf(
		"operator %q does not support %T vs %T",
		o, actual, expected,
	)
}

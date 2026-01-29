package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperator(t *testing.T) {
	type operation struct {
		operator Operator
		a        any
		b        any
	}
	tests := []struct {
		name      string
		operation operation
		expected  bool
		wantErr   require.ErrorAssertionFunc
	}{
		// Numeric comparisons
		{"numeric eq - true", operation{"eq", 25, 25}, true, require.NoError},
		{"numeric eq - false", operation{"eq", 25, 30}, false, require.NoError},
		{"numeric neq - true", operation{"neq", 25, 30}, true, require.NoError},
		{"numeric neq - false", operation{"neq", 25, 25}, false, require.NoError},
		{"numeric lt - true", operation{"lt", 25, 30}, true, require.NoError},
		{"numeric lt - false", operation{"lt", 30, 25}, false, require.NoError},
		{"numeric lte - true (less)", operation{"lte", 25, 30}, true, require.NoError},
		{"numeric lte - true (equal)", operation{"lte", 30, 30}, true, require.NoError},
		{"numeric lte - false", operation{"lte", 31, 30}, false, require.NoError},
		{"numeric gt - true", operation{"gt", 30, 25}, true, require.NoError},
		{"numeric gt - false", operation{"gt", 25, 30}, false, require.NoError},
		{"numeric gte - true (greater)", operation{"gte", 30, 25}, true, require.NoError},
		{"numeric gte - true (equal)", operation{"gte", 25, 25}, true, require.NoError},
		{"numeric gte - false", operation{"gte", 20, 25}, false, require.NoError},

		// String comparisons
		{"string eq - true", operation{"eq", "apple", "apple"}, true, require.NoError},
		{"string eq - false", operation{"eq", "apple", "banana"}, false, require.NoError},
		{"string neq - true", operation{"neq", "apple", "banana"}, true, require.NoError},
		{"string neq - false", operation{"neq", "apple", "apple"}, false, require.NoError},
		{"string lt - true", operation{"lt", "apple", "banana"}, true, require.NoError},
		{"string lt - false", operation{"lt", "banana", "apple"}, false, require.NoError},
		{"string lte - true (less)", operation{"lte", "apple", "banana"}, true, require.NoError},
		{"string lte - true (equal)", operation{"lte", "apple", "apple"}, true, require.NoError},
		{"string lte - false", operation{"lte", "banana", "apple"}, false, require.NoError},
		{"string gt - true", operation{"gt", "banana", "apple"}, true, require.NoError},
		{"string gt - false", operation{"gt", "apple", "banana"}, false, require.NoError},
		{"string gte - true (greater)", operation{"gte", "banana", "apple"}, true, require.NoError},
		{"string gte - true (equal)", operation{"gte", "apple", "apple"}, true, require.NoError},
		{"string gte - false", operation{"gte", "apple", "banana"}, false, require.NoError},

		// Contains
		{"contains string - true", operation{"contains", "gold member", "gold"}, true, require.NoError},
		{"contains string - false", operation{"contains", "gold member", "silver"}, false, require.NoError},
		{"contains slice (string) - true", operation{"contains", "bde", []any{"a", "b", "c"}}, true, require.NoError},
		{"contains slice (string) - false", operation{"contains", "def", []any{"a", "b", "c"}}, false, require.NoError},
		{"contains non-slice/non-string - error", operation{"contains", 123, "a"}, false, require.Error},

		// Like (regexp)
		{"like regexp - true", operation{"like", "user@example.com", ".*@example\\.com$"}, true, require.NoError},
		{"like regexp - false", operation{"like", "user@gmail.com", ".*@example\\.com$"}, false, require.NoError},
		{"like invalid regexp - error", operation{"like", "user@example.com", "["}, false, require.Error},

		// Mixed types / Errors
		{"gt mixed types (int vs string) - error", operation{"gt", 25, "20"}, false, require.Error},
		{"lt mixed types (string vs int) - error", operation{"lt", "25", 20}, false, require.Error},
		{"gt unsupported type (bool) - error", operation{"gt", true, false}, false, require.Error},
		{"unknown operator - error", operation{"unknown", 1, 1}, false, require.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.operation.operator.Evaluate(tt.operation.a, tt.operation.b)
			tt.wantErr(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

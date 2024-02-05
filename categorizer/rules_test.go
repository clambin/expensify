package categorizer_test

import (
	"bytes"
	"github.com/clambin/expensify/categorizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadRules(t *testing.T) {
	input := `
foo:
  - amount: 10

bar:
  - description: ^payment to bar

combo:
  - description: ^payment to foo
    amount: 20

multirule:
  - description: ^payment to foo
  - amount: 20
`

	rules, err := categorizer.LoadRules(bytes.NewBufferString(input))
	assert.NoError(t, err)
	assert.Len(t, rules, 4)

	require.Contains(t, rules, categorizer.Category("foo"))
	require.Len(t, rules["foo"], 1)
	assert.Equal(t, 10.0, rules["foo"][0].Amount)

	require.Contains(t, rules, categorizer.Category("bar"))
	require.Len(t, rules["bar"], 1)
	assert.Equal(t, "^payment to bar", rules["bar"][0].Description)

	require.Contains(t, rules, categorizer.Category("combo"))
	require.Len(t, rules["combo"], 1)
	assert.Equal(t, "^payment to foo", rules["combo"][0].Description)
	assert.Equal(t, 20.0, rules["combo"][0].Amount)

	require.Contains(t, rules, categorizer.Category("multirule"))
	require.Len(t, rules["multirule"], 2)
	assert.Equal(t, "^payment to foo", rules["multirule"][0].Description)
	assert.Equal(t, 20.0, rules["multirule"][1].Amount)
}

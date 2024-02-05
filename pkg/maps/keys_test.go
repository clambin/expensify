package maps_test

import (
	"github.com/clambin/expensify/pkg/maps"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortedKeys(t *testing.T) {
	ints := map[string]int{
		"one": 1,
		"two": 2,
	}

	assert.Equal(t, []string{"one", "two"}, maps.SortedKeys(ints))
}

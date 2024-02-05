package categorizer_test

import (
	"github.com/clambin/expensify/categorizer"
	"github.com/clambin/expensify/payment/checking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestCategorize(t *testing.T) {
	transfers := categorizer.Payments{
		checking.Transfer{
			ValueDate:   time.Date(2022, time.December, 19, 0, 0, 0, 0, time.UTC),
			Amount:      10.0,
			Description: "payment to foo",
		},
		checking.Transfer{
			ValueDate:   time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC),
			Amount:      20.0,
			Description: "payment to bar",
		},
		checking.Transfer{
			ValueDate:   time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC),
			Amount:      -20.0,
			Description: "reimbursement from bar",
		},
	}

	f, err := os.CreateTemp("", "")
	require.NoError(t, err)
	_, _ = f.WriteString(`
payment:
  - description: ^payment to
`)
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	c, err := categorizer.New(f.Name())
	require.NoError(t, err)
	c.Process(transfers)

	matched := c.Matched()
	require.Contains(t, matched, categorizer.Category("payment"))
	assert.Len(t, matched["payment"], 2)

	assert.Len(t, c.Unmatched(), 1)

	c.Reset()
	assert.Empty(t, c.Matched())
	assert.Empty(t, c.Unmatched())
}

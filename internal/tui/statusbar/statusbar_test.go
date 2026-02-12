package statusbar_test

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/clambin/expensify/internal/tui/statusbar"
	"github.com/stretchr/testify/assert"
)

func TestStatusBar(t *testing.T) {
	var s tea.Model = statusbar.New(nil, spinner.WithSpinner(spinner.Dot))

	s, _ = s.Update(s.Init())
	s, _ = s.Update(statusbar.Msg{Text: "test", Warn: true, Spinner: true})

	assert.Equal(t, "test ⣾    ", s.(statusbar.Model).Width(10).View())
}

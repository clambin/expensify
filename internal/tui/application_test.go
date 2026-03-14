package tui

import (
	"bytes"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplication(t *testing.T) {
	app := New(fakeRepo{files: map[string][]byte{
		"file1.csv": []byte(`null,null,null,null,null,null,null
28/01/2026,29/01/2026,"-81,40",EUR,Message 1,Wisselkoers##,Gerelateerde kost##
`),
		"file2.csv": nil,
	}}, nil, DefaultKeyMap())
	tm := teatest.NewTestModel(t, app, teatest.WithInitialTermSize(130, 19))

	// wait for the files to be listed
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file1.csv"))
	}, teatest.WithDuration(time.Second), teatest.WithCheckInterval(10*time.Millisecond))

	// load the first file
	enter := tea.KeyPressMsg{Text: "enter"}
	tm.Send(enter)

	// wait for the summary to be rendered
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Total"))
	})

	// open the first category
	tm.Send(enter)

	// wait for the first category to be rendered
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Message 1"))
	})

	// open the details of the first messafe
	tm.Send(enter)

	// wait for the details to be rendered
	// note: this may be a bit flaky, as "Details" already appears in the list output
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Details"))
	})

	tm.Send(tea.KeyPressMsg{Text: "ctrl+c"})
	tm.WaitFinished(t)

	golden.RequireEqual(t, tm.FinalModel(t).View().Content)
}

func TestApplication_Update_Navigation(t *testing.T) {
	move := func(t testing.TB, app tea.Model, key tea.KeyMsg) tea.Model {
		var cmd tea.Cmd
		app, cmd = app.Update(key)
		require.NotNil(t, cmd)
		msg := cmd()
		require.IsType(t, setActivePaneMsg{}, msg)
		app, _ = app.Update(msg)
		return app
	}

	app := New(fakeRepo{}, nil, DefaultKeyMap())

	// Next pane
	for _, expectedPane := range []paneID{repoPane, summaryPane, statementsPane, repoPane} {
		assert.Equal(t, expectedPane, app.(Application).activePane)
		app = move(t, app, tea.KeyPressMsg{Text: "tab"})
	}

	// Previous pane
	for _, expectedPane := range []paneID{summaryPane, repoPane, statementsPane} {
		assert.Equal(t, expectedPane, app.(Application).activePane)
		app = move(t, app, tea.KeyPressMsg{Text: "shift+tab"})
	}
}

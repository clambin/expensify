package tui

import (
	"errors"
	"maps"
	"slices"
	"testing"

	tea "charm.land/bubbletea/v2"
	"codeberg.org/clambin/bubbles/statusbar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
// TestRepoViewInit tests initialization of the RepoView pane
func TestRepoView_Init(t *testing.T) {
	sb := statusbar.New()
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil, "file2.csv": nil}}
	files, _ := r.List()
	rv := newRepoView(r, nil, defaultRepoKeyMap(), &sb)

	// Init should return a batch that includes a populateRepoFilesMsg
	cmd := rv.Init()
	require.NotNil(t, cmd)
	require.IsType(t, tea.BatchMsg{}, cmd())
	for _, c := range cmd().(tea.BatchMsg) {
		m := c()
		// validate that we loaded the correct files
		if m2, ok := m.(populateRepoFilesMsg); ok {
			assert.Equal(t, files, m2.files)
		}
		rv, _ = rv.Update(m)
	}

	// scroll down and check each file entry
	for i := range files {
		require.Equal(t, files[i], rv.SelectedRow()[0])
		rv, _ = rv.Update(tea.KeyPressMsg{Text: "down"})
	}
}
*/

// TestRepoView_loadStatementsFiles tests loading the files in the repo and populating the table
func TestRepoView_loadRepoFilesCmd(t *testing.T) {
	sb := statusbar.New()
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil, "file2.csv": nil}}
	rv := newRepoView(r, nil, defaultRepoKeyMap(), &sb)

	// load the files
	for _, msg := range unbatch(rv.loadRepoFilesCmd()) {
		rv, _ = rv.Update(msg)
	}

	// the first file should be selected
	row := rv.SelectedRow()
	require.NotNil(t, row)
	assert.Equal(t, "file1.csv", row[0])

	// open the selected file. this should fire off messages to load the file and switch to the summary view
	_, cmd := rv.Update(tea.KeyPressMsg{Text: "enter"})
	require.NotNil(t, cmd)
	msg := cmd()
	require.IsType(t, tea.BatchMsg{}, msg)
	for _, m := range msg.(tea.BatchMsg) {
		msg := m()
		if msg, ok := msg.(populateStatementsMsg); ok {
			assert.Equal(t, "file1.csv", msg.taggedStatements)
		}
	}
}

// TestRepoView_loadRepoFilesCmd_Error validates a failed attempt to load the files
func TestRepoView_loadRepoFilesCmd_Error(t *testing.T) {
	sb := statusbar.New()
	r := fakeRepo{err: assert.AnError}

	rv := newRepoView(r, nil, defaultRepoKeyMap(), &sb)
	cmd := rv.loadRepoFilesCmd()
	require.NotNil(t, cmd)
	sb, _ = sb.Update(cmd())
	assert.Equal(t, "Error loading files: "+assert.AnError.Error(), sb.View())
}

// TestRepoView tests reloading the files
func TestRepoView_reload(t *testing.T) {
	sb := statusbar.New()
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil}}

	rv := newRepoView(r, nil, defaultRepoKeyMap(), &sb)
	_, cmd := rv.Update(tea.KeyPressMsg{Text: "r"})
	msg := cmd()
	require.IsType(t, tea.BatchMsg{}, msg)
	for _, c := range msg.(tea.BatchMsg) {
		if msg, ok := c().(populateRepoFilesMsg); ok {
			assert.Equal(t, []string{"file1.csv"}, msg.files)
		}
	}
}

// TestRepoView_openStatementsFileCmd tests loading a file from the repo
func TestRepoView_openStatementsFileCmd(t *testing.T) {
	sb := statusbar.New()
	r := fakeRepo{files: make(map[string][]byte)}
	rv := newRepoView(&r, nil, defaultRepoKeyMap(), &sb)

	// invalid file
	cmd := rv.openStatementsFileCmd("file1.csv")
	for _, msg := range unbatch(cmd) {
		sb, _ = sb.Update(msg)
	}
	require.Equal(t, "Error loading file file1.csv: read: not found", sb.View())

	// valid file
	r.files["file1.csv"] = []byte(`null,null,null,null,null,null,null
28/01/2026,29/01/2026,"-81,40",EUR,Message 1,Wisselkoers##,Gerelateerde kost##
`)
	cmd = rv.openStatementsFileCmd("file1.csv")
	msgs := unbatch(cmd)
	require.Len(t, msgs, 2)
	require.IsType(t, populateStatementsMsg{}, msgs[1])
	assert.Len(t, msgs[1].(populateStatementsMsg).taggedStatements, 1)
	assert.Equal(t, "bnp-visa", msgs[1].(populateStatementsMsg).file.SchemaName)
}

type fakeRepo struct {
	files map[string][]byte
	err   error
}

func (f fakeRepo) Add(string, []byte) error { return assert.AnError }
func (f fakeRepo) List() ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	files := slices.Collect(maps.Keys(f.files))
	slices.Sort(files)
	return files, nil
}
func (f fakeRepo) Get(key string) ([]byte, error) {
	if f.err != nil {
		return nil, f.err
	}
	content, ok := f.files[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return content, nil
}

func unbatch(cmd tea.Cmd) []tea.Msg {
	msg := cmd()
	b, ok := msg.(tea.BatchMsg)
	if !ok {
		return []tea.Msg{msg}
	}
	msgs := make([]tea.Msg, len(b))
	for i, c := range b {
		msgs[i] = c()
	}
	return msgs
}

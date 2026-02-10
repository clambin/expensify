package tui

import (
	"errors"
	"maps"
	"slices"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRepoViewInit tests initialization of the RepoView pane
func TestRepoView_Init(t *testing.T) {
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil, "file2.csv": nil}}
	files, _ := r.List()
	rv := newRepoView(r, nil, defaultRepoKeyMap())

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
		_ = rv.Update(m)
	}

	// scroll down and check each file entry
	for i := range files {
		require.Equal(t, files[i], rv.SelectedRow[0])
		_ = rv.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
}

// TestRepoView_loadStatementsFiles tests loading the files in the repo and populating the table
func TestRepoView_loadRepoFilesCmd(t *testing.T) {
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil, "file2.csv": nil}}
	rv := newRepoView(r, nil, defaultRepoKeyMap())

	// load the files
	rv.Update(rv.loadRepoFilesCmd()())

	// the first file should be selected
	assert.Equal(t, "file1.csv", rv.SelectedRow[0])

	// open the selected file. this should fire off messages to load the file and switch to the summary view
	cmd := rv.Update(tea.KeyMsg{Type: tea.KeyEnter})
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
	r := fakeRepo{err: assert.AnError}
	rv := newRepoView(r, nil, defaultRepoKeyMap())

	cmd := rv.loadRepoFilesCmd()
	require.NotNil(t, cmd)
	msg := cmd()
	require.IsType(t, statusMsg{}, msg)
	assert.Equal(t, "error loading files: "+assert.AnError.Error(), msg.(statusMsg).text)
}

// TestRepoView tests reloading the files
func TestRepoView_reload(t *testing.T) {
	r := fakeRepo{files: map[string][]byte{"file1.csv": nil}}

	rv := newRepoView(r, nil, defaultRepoKeyMap())
	cmd := rv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
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
	r := fakeRepo{files: make(map[string][]byte)}
	rv := newRepoView(&r, nil, defaultRepoKeyMap())

	// invalid file
	msg := rv.openStatementsFileCmd("file1.csv")()
	require.IsType(t, statusMsg{}, msg)
	require.True(t, msg.(statusMsg).error)

	// valid file
	r.files["file1.csv"] = []byte(`null,null,null,null,null,null,null
28/01/2026,29/01/2026,"-81,40",EUR,Message 1,Wisselkoers##,Gerelateerde kost##
`)
	msg = rv.openStatementsFileCmd("file1.csv")()
	require.IsType(t, populateStatementsMsg{}, msg)
	assert.Len(t, msg.(populateStatementsMsg).taggedStatements, 1)
	assert.Equal(t, "bnp-visa", msg.(populateStatementsMsg).file.SchemaName)
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

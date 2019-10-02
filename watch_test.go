package gowatch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/fsnotify/fsnotify"
)

var (
	assertErrorMsg     = "Expected: %v; Got %v"
	unexpectedErrorMsg = "Unexpected error: %v"
)

type watcherAppTest struct{}

func (wa watcherAppTest) compile() error              { return nil }
func (wa watcherAppTest) start() (*exec.Cmd, error)   { return nil, nil }
func (wa watcherAppTest) restart(cmd *exec.Cmd) error { return errors.New("program should not restart") }

func TestRestartIgnore(t *testing.T) {
	w := watcher{
		dir:    "./testdata/http-server",
		ignore: []string{"main.go"},
		app:    watcherAppTest{},
	}
	event := fsnotify.Event{
		Name: "main.go",
	}

	if err := w.restart(&exec.Cmd{}, event); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
}

func TestIsToIgnoreFile(t *testing.T) {
	w := watcher{
		ignore: []string{"*_test.go"},
	}
	matched, err := w.isToIgnoreFile("main_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("main_test.go should match with pattern *_test.go")
	}
	matched, err = w.isToIgnoreFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Errorf("main.go should not match with pattern *_test.go")
	}

}

func TestDiscoverSubDirectories(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	baseDir := fmt.Sprintf("%s/testdata/http-server", pwd)
	expectedDirectories := []string{baseDir, fmt.Sprintf("%s/foo", baseDir)}
	directories, err := discoverSubDirectories(baseDir)
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if len(directories) != len(expectedDirectories) {
		t.Fatalf(assertErrorMsg, len(expectedDirectories), len(directories))
	}
	for i, dir := range directories {
		if dir != expectedDirectories[i] {
			t.Errorf(assertErrorMsg, expectedDirectories[i], dir)
		}
	}
}

func TestGetCurrentFolderName(t *testing.T) {
	dir := "/home/unittest/gowatch/testcase"
	folderExpected := "testcase"
	folder := getCurrentFolderName(dir)
	if folder != folderExpected {
		t.Fatalf(assertErrorMsg, folderExpected, folder)
	}
}

func TestGetCurrentFolderNameEndWithSlash(t *testing.T) {
	dir := "/home/unittest/gowatch/testcase/"
	folderExpected := "testcase"
	folder := getCurrentFolderName(dir)
	if folder != folderExpected {
		t.Fatalf(assertErrorMsg, folderExpected, folder)
	}
}

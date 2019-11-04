package gowatch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/fsnotify/fsnotify"
)

var (
	assertErrorMsg                 = "Expected: %v; Got %v"
	unexpectedErrorMsg             = "Unexpected error: %v"
	errProgramShoultNotRestartTest = errors.New("program should not restart")
)

type appTestCompileError struct{}

func (wa appTestCompileError) compile() error            { return nil }
func (wa appTestCompileError) start() (*exec.Cmd, error) { return nil, nil }
func (wa appTestCompileError) restart(cmd *exec.Cmd) error {
	return errProgramShoultNotRestartTest
}

type appTest struct{}

func (wa appTest) compile() error              { return nil }
func (wa appTest) start() (*exec.Cmd, error)   { return nil, nil }
func (wa appTest) restart(cmd *exec.Cmd) error { return nil }

func createTmpDir(prefix string) (string, error) {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func createTmpGoFile(dir, pattern string) (*os.File, error) {
	return os.Create(fmt.Sprintf("%s/%s", dir, pattern))
}

func TestAddNewDirectories(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	tmpDir, err := ioutil.TempDir("", "TestAddNewDirectories")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(tmpDir)
	existTmpdir, err := ioutil.TempDir(tmpDir, "TestAddNewDirectories")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	currentDirectories := []string{existTmpdir}
	if _, err = ioutil.TempDir(tmpDir, "TestAddNewDirectories"); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if err := addNewDirectories(w, tmpDir, currentDirectories); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestWriteEventRestart(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := createTmpDir("TestWriteEventRestart")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(dir)
	if err := w.Add(dir); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	watcher := watcher{
		app: appTestCompileError{},
	}

	tmpFile, err := createTmpGoFile(dir, "main.go")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	t.Logf("Create tmp file: %s\n", tmpFile.Name())

	// create file event
	//should pass
	if err := watcher.writeEvent(w, nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	if _, err := tmpFile.Write([]byte(`package main
	func main() {}`)); err != nil {
		t.Fatal(err)
	}

	// write event
	if err := watcher.writeEvent(w, nil); err != nil {
		if !errors.Is(err, errProgramShoultNotRestartTest) {
			t.Errorf(unexpectedErrorMsg, err)
		}
	}

}

func TestWriteEvent(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := createTmpDir("TestWriteEvent")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(dir)
	if err := w.Add(dir); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	watcher := watcher{
		ignore: []string{"/tmp/*/main.go"},
	}

	tmpFile, err := createTmpGoFile(dir, "main.go")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	t.Logf("Create tmp file: %s\n", tmpFile.Name())

	// create file event
	//should pass
	if err := watcher.writeEvent(w, nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	if _, err := tmpFile.Write([]byte(`package main
	func main() {}`)); err != nil {
		t.Fatal(err)
	}

	// write event
	if err := watcher.writeEvent(w, nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestHasNewDirectoriesFalse(t *testing.T) {
	currentDirectories := []string{}
	dir := "./testdata/helloworld/"
	newDirectories, exist, err := hasNewDirectories(dir, currentDirectories)
	if err != nil {
		t.Fatal(unexpectedErrorMsg, err)
	}
	if exist {
		t.Errorf("should not return new directories")
	}
	if len(newDirectories) != 0 {
		t.Errorf("should not return new directories %v\n", newDirectories)
	}
}

func TestHasNewDirectories(t *testing.T) {
	currentDirectories := []string{}
	dir := "./testdata/http-server/"
	newDirectories, exist, err := hasNewDirectories(dir, currentDirectories)
	if err != nil {
		t.Fatal(unexpectedErrorMsg, err)
	}
	if !exist {
		t.Errorf("should return new directories")
	}
	if len(newDirectories) == 0 {
		t.Errorf("should return new directories")
	}
}

func TestRestart(t *testing.T) {
	w := watcher{
		dir: "./testdata/http-server",
		app: appTest{},
	}
	event := fsnotify.Event{
		Name: "main.go",
	}

	if err := w.restart(&exec.Cmd{}, event); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
}

func TestRestartIgnore(t *testing.T) {
	w := watcher{
		dir:    "./testdata/http-server",
		ignore: []string{"main.go"},
		app:    appTestCompileError{},
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

func TestContains(t *testing.T) {
	list := []string{"gowatch"}
	value := list[0]
	if !contains(list, value) {
		t.Errorf("%s exist in %v\n", value, list)
	}
}

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

func (wa appTestCompileError) Compile() error            { return nil }
func (wa appTestCompileError) Start() (*exec.Cmd, error) { return nil, nil }
func (wa appTestCompileError) Stop(cmd *exec.Cmd) error  { return nil }
func (wa appTestCompileError) Restart(cmd *exec.Cmd) error {
	return errProgramShoultNotRestartTest
}

type appTest struct{}

func (wa appTest) Compile() error              { return nil }
func (wa appTest) Start() (*exec.Cmd, error)   { return nil, nil }
func (wa appTest) Stop(cmd *exec.Cmd) error    { return nil }
func (wa appTest) Restart(cmd *exec.Cmd) error { return nil }

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

func TestRun(t *testing.T) {
	w, err := NewWatcher("./testdata/helloworld/", []string{}, []string{}, []string{})
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	errCh := make(chan error)
	go func() {
		errCh <- w.Run()
	}()
	//signal to stop events function
	w.stop <- true
	if err := <-errCh; err != nil {
		if err != ErrStopNotifyEvents {
			t.Fatalf(unexpectedErrorMsg, err)
		}
	}
}

func TestAddNewDirectories(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	watcher := Watcher{watcher: w}

	tmpDir, err := ioutil.TempDir("", "TestAddNewDirectories")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(tmpDir)
	if err := watcher.addDirectories(tmpDir); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}

	newTmpdir, err := ioutil.TempDir(tmpDir, "TestAddNewDirectories")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}

	if err := watcher.addDirectories(newTmpdir); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestEventsRestart(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}

	watcher := Watcher{
		app:     appTestCompileError{},
		watcher: w,
	}

	baseDir, err := createTmpDir("TestWriteEventRestart")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(baseDir)

	if err := watcher.addDirectories(baseDir); err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir(baseDir, "secondDir")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	// create new dir event
	if err := watcher.events(nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}

	defer os.RemoveAll(dir)
	if err := w.Add(dir); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}

	tmpFile, err := createTmpGoFile(dir, "main.go")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	t.Logf("Create tmp file: %s\n", tmpFile.Name())

	// create file event
	if err := watcher.events(nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	if _, err := tmpFile.Write([]byte(`package main
	func main() {}`)); err != nil {
		t.Fatal(err)
	}

	// write event
	if err := watcher.events(nil); err != nil {
		if !errors.Is(err, errProgramShoultNotRestartTest) {
			t.Errorf(unexpectedErrorMsg, err)
		}
	}

}

func TestEvents(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	baseDir, err := createTmpDir("TestWriteEvent")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	defer os.RemoveAll(baseDir)

	watcher := Watcher{
		ignore:  []string{"/tmp/*/*/main.go"},
		watcher: w,
	}

	if err := watcher.addDirectories(baseDir); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	dir, err := ioutil.TempDir(baseDir, "secondDir")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}

	// create new dir event
	if err := watcher.events(nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}

	tmpFile, err := createTmpGoFile(dir, "main.go")
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	t.Logf("Create tmp file: %s\n", tmpFile.Name())

	// create file event
	if err := watcher.events(nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	if _, err := tmpFile.Write([]byte(`package main
	func main() {}`)); err != nil {
		t.Fatal(err)
	}

	// write event
	if err := watcher.events(nil); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestRestart(t *testing.T) {
	w := Watcher{
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
	w := Watcher{
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

func TestShutdowNil(t *testing.T) {
	w := Watcher{}
	if err := w.shutdown(); err == nil {
		t.Errorf("expected %v error", ErrInotifyNil)
	}
}

func TestShutdow(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	w := Watcher{watcher: watcher}
	if err := w.shutdown(); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if err := watcher.Add("/tmp/"); err == nil {
		t.Error("expected error of inotify instance already closed")
	}
}

func TestIsToIgnoreFile(t *testing.T) {
	w := Watcher{
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

func TestContainsFalse(t *testing.T) {
	list := []string{"gowatch"}
	value := "invalid value"
	if contains(list, value) {
		t.Errorf("%s not exist in %v\n", value, list)
	}
}

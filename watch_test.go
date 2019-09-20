package gowatch

import (
	"fmt"
	"os"
	"testing"
)

var (
	assertErrorMsg     = "Expected: %v; Got %v"
	unexpectedErrorMsg = "Unexpected error: %v"
)

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

func TestCompileProgram(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	w := watcher{
		dir:        fmt.Sprintf("%s/testdata/helloworld", pwd),
		binaryName: fmt.Sprintf("%s/testdata/helloworld/helloworld", pwd),
	}
	if err := w.compileProgram(); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
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

func TestCmdRunBinary(t *testing.T) {
	dir := "/home/unittest/gowatch/testcase"
	binaryName := "testcase"
	pathExpected := fmt.Sprintf("./%s", binaryName)
	cmd := cmdRunBinary(dir, binaryName)
	if cmd.Path != pathExpected {
		t.Errorf(assertErrorMsg, pathExpected, cmd.Path)
	}
}

func TestCmdRunBase(t *testing.T) {
	dir := "/home/unittest/gowatch/testcase"
	command := "testcase"
	args := []string{"unittest", "--foobar"}
	argsExpected := []string{"testcase", "unittest", "--foobar"}
	cmd := cmdRunBase(dir, command, args...)
	if cmd.Dir != dir {
		t.Errorf(assertErrorMsg, dir, cmd.Dir)
	}
	if cmd.Path != command {
		t.Errorf(assertErrorMsg, command, cmd.Path)
	}
	if cmd.Stderr != os.Stderr {
		t.Errorf(assertErrorMsg, "os.Stderr", cmd.Stderr)
	}
	if cmd.Stdout != os.Stdout {
		t.Errorf(assertErrorMsg, "os.Stdout", cmd.Stdout)
	}
	if cmd.Stdin != os.Stdin {
		t.Errorf(assertErrorMsg, "os.Stdin", cmd.Stdin)
	}
	if len(cmd.Args) != len(argsExpected) {
		t.Fatalf(assertErrorMsg, len(argsExpected), len(cmd.Args))
	}
	for i, arg := range cmd.Args {
		if arg != argsExpected[i] {
			t.Errorf(assertErrorMsg, argsExpected[i], arg)
		}
	}

}

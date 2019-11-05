package gowatch

import (
	"fmt"
	"os"
	"testing"
)

func TestStartApp(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	w := AppRunner{
		dir:        fmt.Sprintf("%s/testdata/helloworld", pwd),
		binaryName: fmt.Sprintf("%s/testdata/helloworld/helloworld", pwd),
	}
	if err := w.Compile(); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Start(); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestCompileApp(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	w := AppRunner{
		dir:        fmt.Sprintf("%s/testdata/helloworld", pwd),
		binaryName: fmt.Sprintf("%s/testdata/helloworld/helloworld", pwd),
	}
	if err := w.Compile(); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}
func TestCompileAppWithFlags(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	w := AppRunner{
		dir:        fmt.Sprintf("%s/testdata/helloworld", pwd),
		binaryName: fmt.Sprintf("%s/testdata/helloworld/helloworld", pwd),
		buildFlags: []string{"-x", "-v"},
	}
	if err := w.Compile(); err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestRestartApp(t *testing.T) {
	w := AppRunner{
		binaryName: "http-server",
		dir:        "./testdata/http-server",
	}
	if err := w.Compile(); err != nil {
		t.Fatal(err)
	}

	cmd := cmdRunBinary(w.dir, w.binaryName, w.runFlags...)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := w.Restart(cmd); err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
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
func TestCmdRunBinaryPrefix(t *testing.T) {
	dir := "/home/unittest/gowatch/testcase"
	binaryName := "/home/unittest/gowatch/testcase/testcase"
	cmd := cmdRunBinary(dir, binaryName)
	if cmd.Path != binaryName {
		t.Errorf(assertErrorMsg, binaryName, cmd.Path)
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

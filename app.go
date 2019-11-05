package gowatch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

//App manage go apps
type App interface {
	//Compile compile app
	Compile() error

	//Start start app and return cmd executable
	Start() (*exec.Cmd, error)

	//Restart restart a cmd passed in parameter
	Restart(*exec.Cmd) error
}

//AppRunner struct to compile, start
//and restart golang apps
type AppRunner struct {
	// directory to watcher for changes
	dir string

	//flags to use in binary compiled execution
	runFlags []string

	//flags to use in go build command
	buildFlags []string

	//compiled binary name to execute
	binaryName string
}

func (app AppRunner) Compile() error {
	if _, err := os.Stat(app.binaryName); !os.IsNotExist(err) {
		logrus.Debugf("Removing existing binary buildfile %s\n", app.binaryName)
		if err := os.Remove(app.binaryName); err != nil {
			return fmt.Errorf("error to remove existing binary: %v", err)
		}
	}
	buildFlags := []string{"build"}
	buildFlags = append(buildFlags, app.buildFlags...)
	return cmdRunBase(app.dir, "go", buildFlags...).Run()
}

func (app AppRunner) Start() (*exec.Cmd, error) {
	cmd := cmdRunBinary(app.dir, app.binaryName, app.runFlags...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (app AppRunner) Restart(cmd *exec.Cmd) error {
	logrus.Debugf("Killing current execution %d\n", cmd.Process.Pid)
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error to kill exiting process running: %v", err)
	}
	logrus.Debugf("Recompiling...")
	if err := app.Compile(); err != nil {
		return ErrCmdCompile
	}
	*cmd = *cmdRunBinary(app.dir, app.binaryName, app.runFlags...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error to start program: %v", err)
	}

	return nil
}

func cmdRunBinary(dir, binaryName string, args ...string) *exec.Cmd {
	if strings.HasPrefix(binaryName, "/") {
		return cmdRunBase(dir, binaryName, args...)
	}
	return cmdRunBase(dir, fmt.Sprintf("./%s", binaryName), args...)
}

func cmdRunBase(dir, command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Debugf("Command: %v\n", cmd.Args)
	return cmd
}

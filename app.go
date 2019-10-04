package gowatch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type app interface {
	compile() error
	start() (*exec.Cmd, error)
	restart(*exec.Cmd) error
}

type watcherApp struct {
	// directory to watcher for changes
	dir string

	//flags to use in binary compiled execution
	runFlags []string

	//flags to use in go build command
	buildFlags []string

	//compiled binary name to execute
	binaryName string
}

func (wa watcherApp) compile() error {
	if _, err := os.Stat(wa.binaryName); !os.IsNotExist(err) {
		logrus.Debugf("Removing existing binary buildfile %s\n", wa.binaryName)
		if err := os.Remove(wa.binaryName); err != nil {
			return fmt.Errorf("error to remove existing binary: %v", err)
		}
	}
	buildFlags := []string{"build"}
	for _, buildFlag := range wa.buildFlags {
		buildFlags = append(buildFlags, buildFlag)
	}
	return cmdRunBase(wa.dir, "go", buildFlags...).Run()
}
func (wa watcherApp) start() (*exec.Cmd, error) {
	cmd := cmdRunBinary(wa.dir, wa.binaryName, wa.runFlags...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}
func (wa watcherApp) restart(cmd *exec.Cmd) error {
	logrus.Debugf("Killing current execution %d\n", cmd.Process.Pid)
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error to kill exiting process running: %v", err)
	}
	logrus.Debugf("Recompiling...")
	if err := wa.compile(); err != nil {
		return ErrCmdCompile
	}
	*cmd = *cmdRunBinary(wa.dir, wa.binaryName, wa.runFlags...)
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

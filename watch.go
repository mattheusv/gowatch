package gowatch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type watcher struct {
	// directory to watcher for changes
	dir string

	//flags to use in go build command
	buildFlags []string

	//flags to use in binary compiled execution
	runFlags []string

	//compiled binary name to execute
	binaryName string
}

//Start start the watching for changes  in .go files
func Start(dir string, buildFlags, runFlags []string) error {
	w := watcher{
		dir:        dir,
		buildFlags: buildFlags,
		runFlags:   runFlags,
	}
	return w.watch()
}

func (w watcher) watch() error {
	w.binaryName = getCurrentFolderName(w.dir)
	return w.start()
}

func (w watcher) start() error {
	if err := w.compileProgram(); err != nil {
		return err
	}
	return w.startWatch()
}

func (w watcher) compileProgram() error {
	if _, err := os.Stat(w.binaryName); !os.IsNotExist(err) {
		logrus.Debugf("Removing existing binary buildfile %s\n", w.binaryName)
		if err := os.Remove(w.binaryName); err != nil {
			return fmt.Errorf("error to remove existing binary: %v", err)
		}
	}
	buildFlags := []string{"build"}
	for _, buildFlag := range w.buildFlags {
		buildFlags = append(buildFlags, buildFlag)
	}
	return cmdRunBase(w.dir, "go", buildFlags...).Run()
}

func (w watcher) startWatch() error {
	cmd := cmdRunBinary(w.dir, w.binaryName, w.runFlags...)
	if err := cmd.Start(); err != nil {
		return err
	}
	return w.watchFiles(cmd)
}

func (w watcher) watchFiles(cmd *exec.Cmd) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	directories, err := discoverSubDirectories(w.dir)
	if err != nil {
		return err
	}
	for _, dir := range directories {
		if err := watcher.Add(dir); err != nil {
			return err
		}
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				logrus.Errorf("error to get event: %v\n", err)
				os.Exit(5)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if event.Name[len(event.Name)-3:] == ".go" {
					logrus.Infof("Modified file: %s\n", event.Name)
					logrus.Debugf("Killing current execution %d\n", cmd.Process.Pid)
					if err := cmd.Process.Kill(); err != nil {
						return fmt.Errorf("error to kill exiting process running: %v", err)
					}
					logrus.Info("Recompiling...")
					if err := w.compileProgram(); err != nil {
						return fmt.Errorf("could not compile program: %v", err)
					}
					cmd = cmdRunBinary(w.dir, w.binaryName, w.runFlags...)
					if err := cmd.Start(); err != nil {
						return fmt.Errorf("error to start program: %v", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				logrus.Errorf("error watching files changes: %v\n", err)
				os.Exit(5)
			}
		}
	}
}

func discoverSubDirectories(baseDir string) ([]string, error) {
	directories := []string{}
	if err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return directories, nil
}

func getCurrentFolderName(dir string) string {
	folders := strings.Split(dir, "/")
	currentFolder := folders[len(folders)-1]
	if currentFolder == "" {
		return folders[len(folders)-2]
	}
	return currentFolder
}

func cmdRunBinary(dir, binaryName string, args ...string) *exec.Cmd {
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

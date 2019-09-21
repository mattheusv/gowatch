package gowatch

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

var (
	//ErrCmdCompile go build command failed to compile program error
	ErrCmdCompile = errors.New("error to compile program")
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

	// pattern of files to not watch
	ignore []string

	//don't rebuild program after go fmt changes
	skipFmt bool
}

//Start start the watching for changes  in .go files
func Start(dir string, buildFlags, runFlags, ignore []string, skipFmt bool) error {
	w := watcher{
		dir:        dir,
		buildFlags: buildFlags,
		runFlags:   runFlags,
		ignore:     ignore,
		skipFmt:    skipFmt,
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

func (w watcher) isToIgnoreFile(file string) (bool, error) {
	for _, pattern := range w.ignore {
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			return true, err
		}
		if matched {
			return matched, nil
		}
	}
	return false, nil
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
	lastFileChange := ""
	fileChangeFmt := 0
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return fmt.Errorf("error to get event: %v", err)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				lastFileChange = event.Name
				if w.skipFmt {
					if event.Name == lastFileChange {
						// second change because of fmt
						if fileChangeFmt >= 1 {
							logrus.Debug("Skipping fmt changes")
							fileChangeFmt = 0
							// skipt go fmt change and not rebuild and run again
							continue
						}
						fileChangeFmt++
					}
				}
				if event.Name[len(event.Name)-3:] == ".go" {
					if err := w.restart(cmd, event); err != nil {
						if !errors.Is(err, ErrCmdCompile) {
							return err
						}
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return fmt.Errorf("watching files changes: %v", err)
			}
		}
	}
}

func (w watcher) restart(cmd *exec.Cmd, event fsnotify.Event) error {
	ignore, err := w.isToIgnoreFile(event.Name)
	if err != nil {
		return err
	}
	if !ignore {
		logrus.Infof("Modified file: %s\n", event.Name)
		logrus.Debugf("Killing current execution %d\n", cmd.Process.Pid)
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("error to kill exiting process running: %v", err)
		}
		logrus.Info("Recompiling...")
		if err := w.compileProgram(); err != nil {
			return ErrCmdCompile
		}
		*cmd = *cmdRunBinary(w.dir, w.binaryName, w.runFlags...)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error to start program: %v", err)
		}

	}
	return nil
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

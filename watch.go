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

	//ErrInotifyNil nil instance of fsnotify.Watcher
	ErrInotifyNil = errors.New("inotify instance nil")
)

//Watcher struc to watch  to watch for .go file changes
type Watcher struct {
	// directory to watcher for changes
	dir string

	// pattern of files to not watch
	ignore []string

	//interface to start, restart and build the watched program
	app App

	watcher *fsnotify.Watcher
}

//NewWatcher create watcher struct with all values filled
func NewWatcher(dir string, buildFlags, runFlags, ignore []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		ignore:  ignore,
		dir:     dir,
		watcher: watcher,
		app: AppRunner{
			dir:        dir,
			runFlags:   runFlags,
			buildFlags: buildFlags,
			binaryName: getCurrentFolderName(dir),
		},
	}, nil
}

//Run start the watching for changes  in .go files
func (w Watcher) Run() error {
	if err := w.app.Compile(); err != nil {
		return err
	}
	cmd, err := w.app.Start()
	if err != nil {
		return err
	}
	if err := w.start(cmd); err != nil {
		return err
	}
	return w.shutdown()
}

func (w Watcher) shutdown() error {
	if w.watcher == nil {
		return ErrInotifyNil
	}
	return w.watcher.Close()
}

func (w Watcher) isToIgnoreFile(file string) (bool, error) {
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

func (w Watcher) writeEvent(watcher *fsnotify.Watcher, cmd *exec.Cmd) error {
	select {
	case event, ok := <-watcher.Events:
		if !ok {
			return nil
		}
		if event.Op&fsnotify.Write == fsnotify.Write {
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
			return fmt.Errorf("watcher files changes error: %v", err)
		}
	}
	return nil
}

func addNewDirectories(w *fsnotify.Watcher, dir string, currentDirectories []string) error {
	newDirectories, exist, err := hasNewDirectories(dir, currentDirectories)
	if err != nil {
		return err
	}
	if exist {
		logrus.Debugf("find new directories: %v\n", newDirectories)
		for _, dir := range newDirectories {
			if err := w.Add(dir); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w Watcher) start(cmd *exec.Cmd) error {
	directories, err := discoverSubDirectories(w.dir)
	if err != nil {
		return err
	}

	for _, d := range directories {
		if err := w.watcher.Add(d); err != nil {
			return err
		}
	}

	for {
		if err := addNewDirectories(w.watcher, w.dir, directories); err != nil {
			return err
		}
		if err := w.writeEvent(w.watcher, cmd); err != nil {
			return err
		}
	}

}

func (w Watcher) restart(cmd *exec.Cmd, event fsnotify.Event) error {
	ignore, err := w.isToIgnoreFile(event.Name)
	if err != nil {
		return err
	}
	if !ignore {
		logrus.Debugf("Modified file: %s\n", event.Name)
		return w.app.Restart(cmd)
	}
	return nil
}

func contains(list []string, value string) bool {
	for _, n := range list {
		if value == n {
			return true
		}
	}
	return false
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

func hasNewDirectories(dir string, currentDirectories []string) ([]string, bool, error) {
	newDirectories := []string{}
	directories, err := discoverSubDirectories(dir)
	if err != nil {
		return nil, false, err
	}
	for _, d := range directories {
		if d != dir {
			if !contains(currentDirectories, d) {
				newDirectories = append(newDirectories, d)
			}
		}
	}
	if len(newDirectories) != 0 {
		return newDirectories, true, nil
	}
	return nil, false, nil
}

func getCurrentFolderName(dir string) string {
	folders := strings.Split(dir, "/")
	currentFolder := folders[len(folders)-1]
	if currentFolder == "" {
		return folders[len(folders)-2]
	}
	return currentFolder
}

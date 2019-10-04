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

	// pattern of files to not watch
	ignore []string

	//don't rebuild program after go fmt changes
	skipFmt bool

	//interface to start, restart and build the watched program
	app app
}

//Start start the watching for changes  in .go files
func Start(dir string, buildFlags, runFlags, ignore []string, skipFmt bool) error {
	w := watcher{
		ignore:  ignore,
		skipFmt: skipFmt,
		dir:     dir,
		app: watcherApp{
			dir:        dir,
			runFlags:   runFlags,
			buildFlags: buildFlags,
			binaryName: getCurrentFolderName(dir),
		},
	}
	return w.start()
}

func (w watcher) start() error {
	if err := w.app.compile(); err != nil {
		return err
	}
	cmd, err := w.app.start()
	if err != nil {
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

func (w watcher) writeEvent(watcher *fsnotify.Watcher, cmd *exec.Cmd) error {
	select {
	case event, ok := <-watcher.Events:
		if !ok {
			return nil
		}
		if event.Op&fsnotify.Write == fsnotify.Write {
			if w.skipFmt {
				// TODO implement skip fmt changes
				logrus.Warning("--skip-fmt not implemented")
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
			return fmt.Errorf("watcher files changes error: %v", err)
		}
	}
	return nil
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
	addDirectories := func(w *fsnotify.Watcher, d []string) error {
		for _, dir := range d {
			if err := watcher.Add(dir); err != nil {
				return err
			}
		}
		return nil
	}

	if err := addDirectories(watcher, directories); err != nil {
		return err
	}

	for {
		newDirectories, exist, err := hasNewDirectories(w.dir, directories)
		if err != nil {
			return err
		}
		if exist {
			logrus.Debugf("find new directories: %v\n", newDirectories)
			if err := addDirectories(watcher, newDirectories); err != nil {
				return err
			}
		}
		if err := w.writeEvent(watcher, cmd); err != nil {
			return err
		}
	}
}

func (w watcher) restart(cmd *exec.Cmd, event fsnotify.Event) error {
	ignore, err := w.isToIgnoreFile(event.Name)
	if err != nil {
		return err
	}
	if !ignore {
		logrus.Debugf("Modified file: %s\n", event.Name)
		return w.app.restart(cmd)
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

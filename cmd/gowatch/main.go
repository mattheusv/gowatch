package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/msalcantara/gowatch"
	"github.com/msalcantara/gowatch/cmd/gowatch/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	cfg, err := cli(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		os.Exit(2)
	}

	if err := gowatch.Start(cfg.Dir, cfg.Buildflags, cfg.RunFlags, cfg.Ignore); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error to start gowatch"))
		os.Exit(2)
	}
}

func cli(args []string) (config.Config, error) {
	var (
		cfg                              config.Config
		configFileFlag                   string
		buildFlags, runFlags, ignoreFlag string
		runArgs                          []string
		dirFlag                          string
		verboseFlag                      bool
	)

	{
		a := kingpin.New(filepath.Base(os.Args[0]), "watch for .go files changes and rebuild automaticaly")
		//TODO improvement version management
		a.Version("v0.4.0")
		a.HelpFlag.Short('h')

		a.Flag("config", "config file (default .gowatch.yml)").Short('c').Default(".gowatch.yml").StringVar(&configFileFlag)

		a.Flag("build-flags", "flags to go build command").StringVar(&buildFlags)

		a.Flag("run-flags", "custon args to your app").StringVar(&runFlags)

		a.Flag("dir", "directory to wath .go files").Short('d').Default(".").StringVar(&dirFlag)

		a.Flag("ignore", "pattern of files to not watch").Short('i').StringVar(&ignoreFlag)

		a.Flag("verbose", "verbose mode").Short('V').BoolVar(&verboseFlag)

		a.Arg("your-args", "custon args to your app").StringsVar(&runArgs)
		_, err := a.Parse(args)
		if err != nil {
			return config.Config{}, err
		}
	}

	cfg, err := config.LoadYml(configFileFlag)
	if err != nil {
		if !os.IsNotExist(err) {
			return config.Config{}, err
		}
	}
	if len(runArgs) != 0 {
		cfg.RunFlags = runArgs
	}
	if len(runFlags) != 0 {
		cfg.RunFlags = strings.Split(runFlags, ",")
	}
	if len(buildFlags) != 0 {
		cfg.Buildflags = strings.Split(buildFlags, ",")
	}

	if len(ignoreFlag) != 0 {
		cfg.Ignore = strings.Split(ignoreFlag, ",")
	}

	cfg.Dir = dirFlag
	if cfg.Dir == "" || cfg.Dir == "." {
		pwd, err := os.Getwd()
		if err != nil {
			return config.Config{}, err
		}
		cfg.Dir = pwd
	}
	if verboseFlag {
		cfg.Verbose = verboseFlag
	}
	if cfg.Verbose {
		logrus.SetOutput(os.Stdout)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05",
		})
		logrus.SetLevel(logrus.DebugLevel)

	}
	return cfg, nil
}

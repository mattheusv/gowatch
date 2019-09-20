package main

import (
	"fmt"
	"os"

	"github.com/msalcantara/gowatch"
	"github.com/msalcantara/gowatch/cmd/gowatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dir         string
	buildFlags  []string
	runFlags    []string
	ignoreFiles []string
	verbose     bool
)

var rootCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "watching .go files change since 2019-09-13",
	Long:  "gowatch is a tool to watch for .go files changes and rebuild automaticaly",
	Run:   run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "", "directory to wath .go files")
	rootCmd.PersistentFlags().StringSliceVar(&buildFlags, "build-flags", []string{}, "flags to go build command")
	rootCmd.PersistentFlags().StringSliceVar(&runFlags, "run-flags", []string{}, "flags to execute binary")
	rootCmd.PersistentFlags().StringSliceVar(&ignoreFiles, "ignore", []string{}, "pattern of files to not watch")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose mode")
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := initConfig()
	if err != nil {
		exit(err, 3)
	}
	initLogger(cfg.Verbose)
	if err := gowatch.Start(cfg.Dir, cfg.Buildflags, cfg.RunFlags, cfg.Ignore); err != nil {
		exit(err, 3)
	}
}

func initConfig() (config.Config, error) {
	cfg, err := config.LoadYml()
	if err != nil {
		if !os.IsNotExist(err) {
			return config.Config{}, err
		}
	}
	cfg.Dir = dir
	if cfg.Dir == "" || cfg.Dir == "." {
		pwd, err := os.Getwd()
		if err != nil {
			return config.Config{}, err
		}
		cfg.Dir = pwd
	}
	if len(buildFlags) != 0 {
		cfg.Buildflags = buildFlags
	}
	if len(runFlags) != 0 {
		cfg.RunFlags = runFlags
	}
	if verbose != false {
		cfg.Verbose = verbose
	}
	if len(ignoreFiles) != 0 {
		cfg.Ignore = ignoreFiles
	}
	return cfg, nil
}

func initLogger(verbose bool) {
	logrus.SetOutput(os.Stdout)
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func exit(err error, code int) {
	fmt.Printf("ERROR: %v\n", err)
	os.Exit(code)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		exit(err, 2)
	}
}

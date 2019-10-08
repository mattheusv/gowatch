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
	ignoreFiles []string
	verbose     bool
	configFile  string
)

var rootCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "watch for .go files changes and rebuild automaticaly",
	Long:  "gowatch is a tool to watch for .go files changes and rebuild automaticaly",
	Run:   run,
	Example: `
$ gowatch apparg1 apparg2
$ gowatch -d ./custon_dir apparg1 apparg2
$ gowatch -c .gowatch.yml"
$ gowatch --build-flags -x,-v
	`,
	Version: "0.2.0",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "", "directory to wath .go files")
	rootCmd.PersistentFlags().StringSliceVar(&buildFlags, "build-flags", []string{}, "flags to go build command")
	rootCmd.PersistentFlags().StringSliceVarP(&ignoreFiles, "ignore", "i", []string{}, "pattern of files to not watch")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose mode")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".gowatch.yml", "config file")
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := initConfig(args...)
	if err != nil {
		exit(err, 3)
	}
	if cfg.Verbose {
		initLogger()
	}
	if err := gowatch.Start(cfg.Dir, cfg.Buildflags, cfg.RunFlags, cfg.Ignore); err != nil {
		exit(err, 3)
	}
}

func initConfig(args ...string) (config.Config, error) {
	cfg, err := config.LoadYml(configFile)
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
	if len(args) != 0 {
		cfg.RunFlags = args
	}
	if verbose != false {
		cfg.Verbose = verbose
	}
	if len(ignoreFiles) != 0 {
		cfg.Ignore = ignoreFiles
	}
	return cfg, nil
}

func initLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)
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

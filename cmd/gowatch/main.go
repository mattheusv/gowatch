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
	dirFlag         string
	buildFlagsFlag  []string
	ignoreFilesFlag []string
	verboseFlag     bool
	configFileFlag  string
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
	Version: "0.4.0",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dirFlag, "dir", "d", "", "directory to wath .go files")
	rootCmd.PersistentFlags().StringSliceVar(&buildFlagsFlag, "build-flags", []string{}, "flags to go build command")
	rootCmd.PersistentFlags().StringSliceVarP(&ignoreFilesFlag, "ignore", "i", []string{}, "pattern of files to not watch")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "V", false, "verbose mode")
	rootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", ".gowatch.yml", "config file")
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := initConfig(dirFlag, buildFlagsFlag, args, ignoreFilesFlag, verboseFlag)
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

func initConfig(dir string, buildFlags, runFlags, ignoreFlag []string, verbose bool) (config.Config, error) {
	cfg, err := config.LoadYml(dir)
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
	if len(ignoreFlag) != 0 {
		cfg.Ignore = ignoreFlag
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

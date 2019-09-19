package main

import (
	"fmt"
	"os"

	"github.com/msalcantara/gowatch"
	"github.com/spf13/cobra"
)

var (
	dir        string
	buildFlags []string
	runFlags   []string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "gowatch",
	Short: "watching .go files change since 2019-09-13",
	Long:  "gowatch is a tool to watch for .go files changes and rebuild automaticaly",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				os.Exit(3)
			}
			dir = pwd
		}
		if err := gowatch.Start(dir, buildFlags, runFlags, verbose); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(3)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "", "directory to wath .go files")
	rootCmd.PersistentFlags().StringSliceVar(&buildFlags, "build-flags", []string{}, "flags to go build command")
	rootCmd.PersistentFlags().StringSliceVar(&runFlags, "run-flags", []string{}, "flags to execute binary")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose mode")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
		fmt.Printf("ERROR: %v\n", err)
	}
}

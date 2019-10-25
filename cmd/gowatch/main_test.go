package main

import (
	"errors"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

var (
	assertErrorMsg     = "Expected: %v; Got %v"
	unexpectedErrorMsg = "Unexpected error: %v"
)

func contains(list []string, value string) bool {
	for _, n := range list {
		if value == n {
			return true
		}
	}
	return false
}

func TestRunFlagsFlagsArgs(t *testing.T) {
	args := []string{"--run-flags", "localhost,8080-v", "foo", "bar"}
	cfg, err := cli(args)
	if err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	for _, arg := range cfg.RunFlags {
		if !contains(args, arg) {
			t.Errorf("invalid parse: %s not exist in %v[%v]", arg, args, len(cfg.RunFlags))
		}
	}

}

func TestRunFlagsFlags(t *testing.T) {
	args := []string{"--run-flags", "localhost,8080-v"}
	cfg, err := cli(args)
	if err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}
	for _, arg := range cfg.RunFlags {
		if !contains(args, arg) {
			t.Errorf("invalid parse: %s not exist in %v[%v]", arg, args, len(cfg.RunFlags))
		}
	}

}
func TestRunFlagsArgs(t *testing.T) {
	args := []string{"localhost", "8080"}
	cfg, err := cli(args)
	if err != nil {
		t.Errorf(unexpectedErrorMsg, err)
	}

	for _, arg := range cfg.RunFlags {
		if !contains(args, arg) {
			t.Errorf("invalid parse: %s not exist", arg)
		}
	}

}

func TestInitConfigErrorYml(t *testing.T) {
	_, err := cli([]string{"-c", "./testdata/gowatch.yml.invalid"})
	var typeError *yaml.TypeError
	if !errors.As(err, &typeError) {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestInitConfigDirPwd(t *testing.T) {
	cfg, err := cli([]string{})
	if err != nil {
		t.Fatal(unexpectedErrorMsg, err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(unexpectedErrorMsg, err)
	}
	if cfg.Dir != pwd {
		t.Errorf(assertErrorMsg, pwd, cfg.Dir)
	}
}

func TestInitConfigCmdFlags(t *testing.T) {
	errTemplate := "%s don't load correctlly from command line: %v"
	dir := "/tmp/whatever/dir"
	cfg, err := cli([]string{"-d", dir, "-V", "-i", "*_test.go", "--build-flags", "x,v", "localhost 8000"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Verbose {
		t.Errorf(errTemplate, "verbose", cfg.Verbose)
	}
	if cfg.Dir != dir {
		t.Errorf(errTemplate, "dir", cfg.Dir)
	}
	if len(cfg.Ignore) == 0 {
		t.Errorf(errTemplate, "ignore", cfg.Ignore)
	}
	if len(cfg.Buildflags) == 0 {
		t.Errorf(errTemplate, "build-flags", cfg.Buildflags)
	}
	if len(cfg.RunFlags) == 0 {
		t.Errorf(errTemplate, "run-flags", cfg.RunFlags)
	}
}

func TestInitConfigConfigFile(t *testing.T) {
	errTemplate := "%s don't load correctlly from command line: %v"
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(unexpectedErrorMsg, err)
	}

	cfg, err := cli([]string{"-c", "./testdata/gowatch.yml"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Verbose {
		t.Errorf(errTemplate, "verbose", cfg.Verbose)
	}
	if cfg.Dir != pwd {
		t.Errorf(errTemplate, "dir", cfg.Dir)
	}
	if len(cfg.Ignore) == 0 {
		t.Errorf(errTemplate, "ignore", cfg.Ignore)
	}
	if len(cfg.Buildflags) == 0 {
		t.Errorf(errTemplate, "build-flags", cfg.Buildflags)
	}
	if len(cfg.RunFlags) == 0 {
		t.Errorf(errTemplate, "run-flags", cfg.RunFlags)
	}
}

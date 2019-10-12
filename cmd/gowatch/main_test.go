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

func TestInitConfigErrorYml(t *testing.T) {
	_, err := initConfig("./testdata/gowatch.yml.invalid", []string{}, []string{}, []string{}, false)
	var typeError *yaml.TypeError
	if !errors.As(err, &typeError) {
		t.Errorf(unexpectedErrorMsg, err)
	}
}

func TestInitConfigDirPwd(t *testing.T) {
	cfg, err := initConfig("", []string{}, []string{}, []string{}, false)
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
	dir := "/tmp/whatever.yml"
	cfg, err := initConfig(dir, []string{
		"x",
		"v",
	}, []string{
		"8080",
	}, []string{
		"*_test.go",
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Verbose {
		t.Errorf("verbose field don't load correctlly from cmd file")
	}
	if cfg.Dir != dir {
		t.Errorf("dir field don't load correctlly from cmd file")
	}
	if len(cfg.Ignore) == 0 {
		t.Errorf("ignore field don't load correctlly from cmd file")
	}
	if len(cfg.Buildflags) == 0 {
		t.Errorf("build_flags field don't load correctlly from cmd file")
	}
	if len(cfg.RunFlags) == 0 {
		t.Errorf("run_flags field don't load correctlly from cmd file")
	}
}

func TestInitConfigConfigFile(t *testing.T) {
	dir := "./testdata/gowatch.yml"
	cfg, err := initConfig(dir, []string{}, []string{}, []string{}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Verbose {
		t.Errorf("verbose field don't load correctlly from config file")
	}
	if cfg.Dir != dir {
		t.Errorf("dir field don't load correctlly from config file")
	}
	if len(cfg.Ignore) == 0 {
		t.Errorf("ignore field don't load correctlly from config file")
	}
	if len(cfg.Buildflags) == 0 {
		t.Errorf("build_flags field don't load correctlly from config file")
	}
	if len(cfg.RunFlags) == 0 {
		t.Errorf("run_flags field don't load correctlly from config file")
	}
}

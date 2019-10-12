package main

import (
	"errors"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestInitConfigErrorYml(t *testing.T) {
	_, err := initConfig("./testdata/gowatch.yml.invalid", []string{}, []string{}, []string{}, false)
	var typeError *yaml.TypeError
	if !errors.As(err, &typeError) {
		t.Errorf(err.Error())
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

package config

import (
	"testing"
)

var (
	unexpectedErrorMsg        = "Unexpected error: %v"
	keyDontLoadedCorrectlyMsg = "key %s don't loaded correctly"
)

func TestLoadYml(t *testing.T) {
	ymlFile := "../testdata/gowatch.yml"
	cfg, err := LoadYml(ymlFile)
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}

	if cfg.Verbose == false {
		t.Errorf(keyDontLoadedCorrectlyMsg, "verbose")
	}
	if cfg.Dir != "." {
		t.Errorf(keyDontLoadedCorrectlyMsg, "dir")
	}

	if len(cfg.Ignore) != 1 {
		t.Errorf(keyDontLoadedCorrectlyMsg, "ignore")
	}

	if len(cfg.Buildflags) != 2 {
		t.Errorf(keyDontLoadedCorrectlyMsg, "build_flags")
	}

	if len(cfg.RunFlags) != 2 {
		t.Errorf(keyDontLoadedCorrectlyMsg, "run_flags")
	}
}

func TestLoadYmlConfig(t *testing.T) {
	ymlFile := "../testdata/gowatch.yml"
	var config Config
	err := loadYmlConfig(&config, ymlFile)
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if config.Dir == "" {
		t.Fatalf(keyDontLoadedCorrectlyMsg, "dir")
	}
	if len(config.Buildflags) == 0 {
		t.Fatalf(keyDontLoadedCorrectlyMsg, "build_flags")
	}
	if len(config.RunFlags) == 0 {
		t.Fatalf(keyDontLoadedCorrectlyMsg, "run_flags")
	}
	if len(config.Ignore) == 0 {
		t.Fatalf(keyDontLoadedCorrectlyMsg, "ignore")
	}
	if config.Verbose == false {
		t.Fatalf(keyDontLoadedCorrectlyMsg, "verbose")
	}
}

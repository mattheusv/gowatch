package config

import "testing"

var (
	assertErrorMsg     = "Expected: %v; Got %v"
	unexpectedErrorMsg = "Unexpected error: %v"
)

func TestLoadYmlConfig(t *testing.T) {
	ymlFile := "../testdata/gowatch.yml"
	var config Config
	err := loadYmlConfig(&config, ymlFile)
	if err != nil {
		t.Fatalf(unexpectedErrorMsg, err)
	}
	if config.Dir == "" {
		t.Fatalf("key dir don't loaded correctly")
	}
	if len(config.Buildflags) == 0 {
		t.Fatalf("key build_flags don't loaded correctly")
	}
	if len(config.RunFlags) == 0 {
		t.Fatalf("key run_flags don't loaded correctly")
	}
	if len(config.Ignore) == 0 {
		t.Fatalf("key ignore don't loaded correctly")
	}
	if config.Verbose == false {
		t.Fatalf("key verbose don't loaded correctly")
	}
}

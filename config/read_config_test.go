package config

import (
	"os"
	"testing"
)

func TestProcessNoConfigFilePresent(t *testing.T) {

	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without any valid config file should fail.")
	} else {
		if err.Error() != "Fatal error reading config file: Config File \"config\" Not Found in \"[/etc/music-manager]\"" {
			t.Errorf("Default config should be in /etc/music-manager/config.toml, not in other place, error was '%s'.", err.Error())
		}
	}
}

func TestProcessServerNoDataInConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_no_data/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without server data config should fail.")
	} else {
		if err.Error() != "Fatal error config: no server host was found." {
			t.Errorf("Error should be \"Fatal error config: no server host was found.\" but error was '%s'.", err.Error())
		}
	}
}

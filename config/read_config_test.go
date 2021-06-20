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
		if err.Error() != "Fatal error reading config: no server host was found." {
			t.Errorf("Error should be \"Fatal error reading config: no server host was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessServerNoWrappersNeitherServices(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_no_wrappers_neither_services/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without wrappers data config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: no wrappers config was found." {
			t.Errorf("Error should be \"Fatal error reading config: no wrappers config was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessServerNoStatusService(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/no_status_service/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without status service data config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: no status config was found." {
			t.Errorf("Error should be \"Fatal error reading config: no status config was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessServerNoStorageService(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/no_storage_service/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without status service data config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: no storage config was found." {
			t.Errorf("Error should be \"Fatal error reading config: no storage config was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessServerNoWrappersDefined(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_no_wrappers_defined/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without wrappers data config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: no wrappers were found, at least one wrapper must be defined." {
			t.Errorf("Error should be \"Fatal error reading config: no wrappers were found, at least one wrapper must be defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWrapperInvalid(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_one_wrapper_invalid_config/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with wrapper invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: wrapper firstwrapper has an invalid config: durable is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: wrapper firstwrapper has an invalid config: durable is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWrapperInvalid2(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_one_wrapper_invalid_config2/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with wrapper invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: wrapper firstwrapper has an invalid config: auto_ack is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: wrapper firstwrapper has an invalid config: auto_ack is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWrapperValidButSecondInvalid(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/server_one_wrapper_second_invalid/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with wrapper invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: wrapper secondwrapper has an invalid config: exclusive is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: wrapper secondwrapper has an invalid config: exclusive is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWithoutJobManagerConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/no_jobmanager_service/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without jobmanager config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: no jobmanager config was found." {
			t.Errorf("Error should be \"Fatal error reading config: no jobmanager config was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWithInvalidJobManagerConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/invalid_jobmanager_config/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with jobmanager invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: jobmanager has an invalid config: exclusive is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: jobmanager has an invalid config: exclusive is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWithInvalidStatusConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/invalid_status_config/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with status invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: status has an invalid config: name is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: status has an invalid config: name is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessWithInvalidStorageConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/invalid_storage_config/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with storage invalid config should fail.")
	} else {
		if err.Error() != "Fatal error reading config: storage has an invalid config: name is not defined." {
			t.Errorf("Error should be \"Fatal error reading config: storage has an invalid config: name is not defined.\" but error was '%s'.", err.Error())
		}
	}
}

func TestValisConfig(t *testing.T) {
	os.Setenv("MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION", "./config_files_test/valid_config/")
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig method with valid config shouldn't fail.")
	}
}

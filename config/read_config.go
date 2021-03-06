package config

import (
	"errors"

	viperLib "github.com/spf13/viper"
)

type Server struct {
	User     string
	Password string
	Host     string
	Port     int
}

type Queue struct {
	Name string
}

type Config struct {
	Server        Server
	Wrappers      []Queue
	Status        string
	Storage       string
	JobManager    Queue
	WrapperOutput Queue
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION"

	serverVariables := []string{"host", "port", "user", "password"}
	queueVariables := []string{"name"}

	requiredConfigEntities := []string{"wrappers", "status", "storage", "jobmanager", "wrapperoutput"}

	viper := viperLib.New()

	//Look for config file location defined as env var
	viper.BindEnv(envVariable)
	configFileLocation = viper.GetString(envVariable)
	if configFileLocation == "" {
		// Get config file from default location
		configFileLocation = "/etc/music-manager/"
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configFileLocation)

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.New(errors.New("Fatal error reading config file: ").Error() + err.Error())
	}

	for _, server_variable := range serverVariables {
		if !viper.IsSet("server." + server_variable) {
			return config, errors.New("Fatal error reading config: no server " + server_variable + " was found.")
		}
	}

	server := Server{User: viper.GetString("server.user"), Password: viper.GetString("server.password"), Host: viper.GetString("server.host"), Port: viper.GetInt("server.port")}

	config.Server = server

	for _, requiredConfigEntity := range requiredConfigEntities {
		if !viper.IsSet(requiredConfigEntity) {
			return config, errors.New("Fatal error reading config: no " + requiredConfigEntity + " config was found.")
		}
	}

	// Check Wrappers

	wrapperConfig := viper.Get("wrappers")
	wrapperConfigElementsMap := wrapperConfig.(map[string]interface{})

	if len(wrapperConfigElementsMap) == 0 {
		return config, errors.New("Fatal error reading config: no wrappers were found, at least one wrapper must be defined.")
	}

	for wrapperName, _ := range wrapperConfigElementsMap {
		for _, requiredQueueVeriable := range queueVariables {
			if !viper.IsSet("wrappers." + wrapperName + "." + requiredQueueVeriable) {
				return config, errors.New("Fatal error reading config: wrapper " + wrapperName + " has an invalid config: " + requiredQueueVeriable + " is not defined.")
			}
		}
		wrapperConfig := Queue{Name: viper.GetString("wrappers." + wrapperName + ".name")}
		config.Wrappers = append(config.Wrappers, wrapperConfig)
	}

	// Check JobManager
	for _, requiredQueueVeriable := range queueVariables {
		if !viper.IsSet("jobmanager." + requiredQueueVeriable) {
			return config, errors.New("Fatal error reading config: jobmanager has an invalid config: " + requiredQueueVeriable + " is not defined.")
		}
	}
	jobmanagerConfig := Queue{Name: viper.GetString("jobmanager.name")}
	config.JobManager = jobmanagerConfig

	// Check WrapperOutput
	for _, requiredQueueVeriable := range queueVariables {
		if !viper.IsSet("wrapperoutput." + requiredQueueVeriable) {
			return config, errors.New("Fatal error reading config: wrapperoutput has an invalid config: " + requiredQueueVeriable + " is not defined.")
		}
	}
	wrapperoutputConfig := Queue{Name: viper.GetString("wrapperoutput.name")}
	config.WrapperOutput = wrapperoutputConfig

	// Check Status
	if !viper.IsSet("status.name") {
		return config, errors.New("Fatal error reading config: status has an invalid config: name is not defined.")
	}
	config.Status = viper.GetString("status.name")

	// Check Storage
	if !viper.IsSet("storage.name") {
		return config, errors.New("Fatal error reading config: storage has an invalid config: name is not defined.")
	}
	config.Storage = viper.GetString("storage.name")

	return config, nil
}

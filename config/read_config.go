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
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	NoLocal          bool
	AutoACK          bool
}

type Config struct {
	Server     Server
	Wrappers   []Queue
	Status     Queue
	Storage    Queue
	JobManager Queue
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION"

	serverVariables := []string{"host", "port", "user", "password"}
	queueVariables := []string{"name", "durable", "delete_when_unused", "exclusive", "no_wait", "auto_ack"}

	requiredConfigEntities := []string{"wrappers", "status", "storage", "jobs"}

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
	}

	// Check JobManager
	for _, requiredQueueVeriable := range queueVariables {
		if !viper.IsSet("jobs." + requiredQueueVeriable) {
			return config, errors.New("Fatal error reading config: jobs has an invalid config: " + requiredQueueVeriable + " is not defined.")
		}
	}

	// Check Status
	if !viper.IsSet("status.name") {
		return config, errors.New("Fatal error reading config: status has an invalid config: name is not defined.")
	}
	// Check Storage
	if !viper.IsSet("storage.name") {
		return config, errors.New("Fatal error reading config: storage has an invalid config: name is not defined.")
	}

	return config, nil
}

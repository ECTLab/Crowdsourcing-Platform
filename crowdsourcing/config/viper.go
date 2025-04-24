package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	stdlog "log"
	"os"
	"sync"
)

var (
	configPath     string
	once           sync.Once
	configInstance *Config
	loadErr        error
)

func InitViperWithConfigPathForEnv(path, env string) {
	configPath = path + strings.ToLower(env) + ".yaml"
	once = sync.Once{}
}
func InitViper() {
	node_env := os.Getenv("NODE_ENV")
	if node_env == "" {
		node_env = "DEVELOPMENT"
	}

	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		config_path = "config/"
	}

	InitViperWithConfigPathForEnv(config_path, node_env)
}

func newConfig(configPath string) (*Config, error) {
	config, err := newConfigByViper(configPath)
	return config, err
}

func loadConfig() {
	configInstance, loadErr = newConfig(configPath)
	if loadErr != nil {
		stdlog.Printf("Couldn't get config from %s with error: %s", configPath, loadErr)
	} else {
		stdlog.Printf("Successfully get config from %s", configPath)
	}
}

var GetServiceConfig = func() Config {
	if configPath == "" {
		panic(errors.New("using GetServiceConfig() before initializing. call InitViper()"))
	}
	once.Do(loadConfig)

	if loadErr != nil {
		stdlog.Printf("unable to load config file")
	}

	return *configInstance
}

func newConfigByViper(configPath string) (*Config, error) {
	initializeViper(configPath)
	err := readConfigFileByViper()
	if err != nil {
		return nil, err
	}
	config, err := unmarshallFileToConfigByViper()
	if err != nil {
		return nil, err
	}
	config, err = overwriteEnvValuesToConfigByViper(config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func overwriteEnvValuesToConfigByViper(config Config) (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}

func unmarshallFileToConfigByViper() (Config, error) {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func readConfigFileByViper() error {
	err := viper.ReadInConfig()
	return err
}

func initializeViper(configPath string) {
	viper.SetConfigName(configPath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}

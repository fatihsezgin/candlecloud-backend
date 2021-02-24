package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	configuration     *Configuration
	configFileName    = "config"
	configFileExt     = ".yml"
	configType        = "yaml"
	storeDirectory    = "./store/"
	configFileAbsPath = filepath.Join(storeDirectory, configFileName)
)

type Configuration struct {
	Database DatabaseConfiguration
	Server   ServerConfiguration
}

type DatabaseConfiguration struct {
	Name     string `default:"candlecloud"`
	Username string `default:"user"`
	Password string `default:"password"`
	Host     string `default:"localhost"`
	Port     string `default:"5432"`
	LogMode  bool   `default:"false"`
}

type ServerConfiguration struct {
	Env        string `default:"dev"` // dev, prod
	Port       string `default:"3625"`
	Timeout    int    `default:"24"`
	Passphrase string `default:"passphrase-for-encrypting-passwords-do-not-forget"`
}

func SetupConfigDefaults() (*Configuration, error) {

	//initialize viper configuration
	initializeConfig()

	//set default values
	setDefaults()

	if err := readConfiguration(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}

// read configuration from file
func readConfiguration() error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// if file does not exist, simply create one
		if _, err := os.Stat(configFileAbsPath + configFileExt); os.IsNotExist(err) {
			os.Create(configFileAbsPath + configFileExt)
		} else {
			return err
		}
		// let's write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

// initialize the configuration manager
func initializeConfig() {
	viper.AddConfigPath(storeDirectory)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)
}

func setDefaults() {

	// Server defaults
	viper.SetDefault("server.env", "prod")
	viper.SetDefault("server.port", "3625")
	viper.SetDefault("server.timeout", 24)
	viper.SetDefault("server.passphrase", generateKey())

	// Database defaults
	viper.SetDefault("database.name", "candlecloud")
	viper.SetDefault("database.username", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.logmode", false)
}

func generateKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "add-your-key-to-here"
	}
	keyEnc := base64.StdEncoding.EncodeToString(key)
	return keyEnc
}

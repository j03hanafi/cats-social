package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig() (*Runtime, error) {
	runtimeViper := viper.New()

	// Get config from env vars
	readEnv(runtimeViper)

	// Get config from file
	err := readConfigFile(runtimeViper, "config", "toml", "./configs")
	if err != nil {
		return nil, err
	}

	// Load config into runtimeConfig
	runtimeConfig, err := loadConfig(runtimeViper)
	if err != nil {
		return nil, err
	}

	return runtimeConfig, nil
}

func readEnv(runtimeViper *viper.Viper) {
	// Set defaults for env vars
	runtimeViper.SetDefault(dbName, "cats_social")
	runtimeViper.SetDefault(dbPort, 5432)
	runtimeViper.SetDefault(dbHost, "localhost")
	runtimeViper.SetDefault(dbUsername, "postgres")
	runtimeViper.SetDefault(dbPassword, "password")
	runtimeViper.SetDefault(dbParams, []string{"sslmode=disable"})
	runtimeViper.SetDefault(jwtSecret, "secret")
	runtimeViper.SetDefault(bcryptSalt, 8)

	// Load env vars
	runtimeViper.AllowEmptyEnv(false)
	runtimeViper.AutomaticEnv()
}

func readConfigFile(runtimeViper *viper.Viper, fileName, fileType string, filePath ...string) error {
	runtimeViper.SetConfigName(fileName)
	runtimeViper.SetConfigType(fileType)
	for _, path := range filePath {
		runtimeViper.AddConfigPath(path)
	}

	err := runtimeViper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(runtimeViper *viper.Viper) (*Runtime, error) {
	// load env vars to dbCfg and apiCfg
	dbConfig, apiConfig := &dbCfg{}, &apiCfg{}
	err := runtimeViper.Unmarshal(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dbConfig: %v\n", err)
	}
	err = runtimeViper.Unmarshal(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dbConfig: %v\n", err)
	}

	// set env vars to runtimeConfig before decode from config file
	runtimeConfig := &Runtime{
		API: *apiConfig,
		DB:  *dbConfig,
	}
	err = runtimeViper.Unmarshal(runtimeConfig)
	if err != nil {
		return nil, err
	}

	return runtimeConfig, nil
}

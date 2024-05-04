package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

var Runtime *RuntimeConfig

func NewConfig() error {
	callerInfo := "[configs.NewConfig]"

	runtimeViper := viper.New()

	// ListCats config from env vars
	readEnv(runtimeViper)

	// ListCats config from file
	err := readConfigFile(runtimeViper, "config", "toml", "./configs")
	if err != nil {
		return fmt.Errorf("%s failed to read config file: %v\n", callerInfo, err)
	}

	// Load config into runtimeConfig
	Runtime, err = loadConfig(runtimeViper)
	if err != nil {
		return fmt.Errorf("%s failed to load config: %v\n", callerInfo, err)
	}

	return nil
}

func readEnv(runtimeViper *viper.Viper) {
	// Set defaults for env vars
	runtimeViper.SetDefault(dbName, "cats_social")
	runtimeViper.SetDefault(dbPort, 5432)
	runtimeViper.SetDefault(dbHost, "localhost")
	runtimeViper.SetDefault(dbUsername, "cats_social")
	runtimeViper.SetDefault(dbPassword, "password")
	runtimeViper.SetDefault(dbParams, []string{"sslmode=disable"})
	runtimeViper.SetDefault(jwtSecret, "secret")
	runtimeViper.SetDefault(bcryptSalt, 8)

	// Load env vars
	runtimeViper.AllowEmptyEnv(false)
	runtimeViper.AutomaticEnv()
}

func readConfigFile(runtimeViper *viper.Viper, fileName, fileType string, filePath ...string) error {
	callerInfo := "[configs.readConfigFile]"

	runtimeViper.SetConfigName(fileName)
	runtimeViper.SetConfigType(fileType)
	for _, path := range filePath {
		runtimeViper.AddConfigPath(path)
	}

	err := runtimeViper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("%s failed to read config file: %v\n", callerInfo, err)
	}
	return nil
}

func loadConfig(runtimeViper *viper.Viper) (*RuntimeConfig, error) {
	callerInfo := "[configs.loadConfig]"

	// load env vars to dbCfg and apiCfg
	dbConfig, apiConfig := &dbCfg{}, &apiCfg{}
	err := runtimeViper.Unmarshal(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode dbConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode apiConfig: %v\n", callerInfo, err)
	}

	// because we get jwtSecret from env vars but in config runtime it's nested, we need to manually set it
	apiConfig.JWT.JWTSecret = runtimeViper.GetString(jwtSecret)

	// set env vars to runtimeConfig before decode from config file
	runtimeConfig := &RuntimeConfig{
		API: *apiConfig,
		DB:  *dbConfig,
	}
	err = runtimeViper.Unmarshal(runtimeConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode runtimeConfig: %v\n", callerInfo, err)
	}

	return runtimeConfig, nil
}

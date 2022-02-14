package localConfig

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/lithammer/shortuuid"
	"github.com/spf13/viper"
)

type ConfigContent struct {
	CliId         string
	SchemaVersion string
}

type LocalConfig struct {
}

func NewLocalConfig() *LocalConfig {
	return &LocalConfig{}
}

func (lc *LocalConfig) GetLocalConfiguration() (*ConfigContent, error) {
	viper.SetEnvPrefix("datree")
	viper.AutomaticEnv()
	token := viper.GetString("token")
	schemaVersion := viper.GetString("schema_version")

	if token == "" {
		configHome, configName, configType, err := setViperConfig()
		if err != nil {
			return nil, err
		}

		// workaround for creating config file when not exist
		// open issue in viper: https://github.com/spf13/viper/issues/430
		// should be fixed in pr https://github.com/spf13/viper/pull/936
		configPath := filepath.Join(configHome, configName+"."+configType)
		_, err = os.Stat(configPath)

		if err != nil {
			os.Mkdir(configHome, os.ModePerm)
			os.Create(configPath)
			viper.SetDefault("token", shortuuid.New())
			viper.WriteConfig()
		}

		viper.ReadInConfig()
		token = viper.GetString("token")

		if token == "" {
			viper.SetDefault("token", shortuuid.New())
			viper.WriteConfig()
			viper.ReadInConfig()
			token = viper.GetString("token")
		}
	}

	return &ConfigContent{CliId: token, SchemaVersion: schemaVersion}, nil
}

func (lc *LocalConfig) Set(key string, value string) error {
	_, _, _, err := setViperConfig()
	if err != nil {
		return err
	}

	viper.Set(key, value)
	viper.WriteConfig()
	return nil
}

func getConfigHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	homedir := usr.HomeDir
	configHome := filepath.Join(homedir, ".datree")

	return configHome, nil
}

func getConfigName() string {
	return "config"
}

func getConfigType() string {
	return "yaml"
}

func setViperConfig() (string, string, string, error) {
	configHome, err := getConfigHome()
	if err != nil {
		return "", "", "", nil
	}

	configName := getConfigName()
	configType := getConfigType()

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configHome)

	return configHome, configName, configType, nil
}

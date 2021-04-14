package localConfig

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/lithammer/shortuuid"
	"github.com/spf13/viper"
)

type LocalConfiguration struct {
	CliId string
}

func (l *LocalConfiguration) GetConfiguration() (LocalConfiguration, error) {
	viper.SetEnvPrefix("datree")
	viper.AutomaticEnv()
	token := viper.GetString("token")

	if token == "" {
		usr, err := user.Current()
		if err != nil {
			return LocalConfiguration{}, err
		}

		homedir := usr.HomeDir

		configHome := filepath.Join(homedir, ".datree")
		configName := "config"
		configType := "yaml"

		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.AddConfigPath(configHome)

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

	return LocalConfiguration{CliId: token}, nil
}

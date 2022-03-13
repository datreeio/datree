package localConfig

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/lithammer/shortuuid"
	"github.com/spf13/viper"
)

type LocalConfig struct {
	Token         string
	ClientId      string
	SchemaVersion string
}

type TokenClient interface {
	CreateToken() (*cliClient.CreateTokenResponse, error)
}

type LocalConfigClient struct {
	tokenClient TokenClient
}

func NewLocalConfigClient(t TokenClient) *LocalConfigClient {
	return &LocalConfigClient{
		tokenClient: t,
	}
}

const (
	clientIdKey      = "client_id"
	tokenKey         = "token"
	schemaVersionKey = "schema_version"
)

func (lc *LocalConfigClient) GetLocalConfiguration() (*LocalConfig, error) {
	viper.SetEnvPrefix("datree")
	viper.AutomaticEnv()

	initConfigFileErr := InitLocalConfigFile()
	if initConfigFileErr != nil {
		return nil, initConfigFileErr
	}

	token := viper.GetString(tokenKey)
	clientId := viper.GetString(clientIdKey)
	schemaVersion := viper.GetString(schemaVersionKey)

	if token == "" {
		createTokenResponse, err := lc.tokenClient.CreateToken()
		if err != nil {
			return nil, err
		}
		token = createTokenResponse.Token
		viper.SetDefault(tokenKey, token)
		writeTokenErr := viper.WriteConfig()
		if writeTokenErr != nil {
			return nil, writeTokenErr
		}
		readTokenErr := viper.ReadInConfig()
		if readTokenErr != nil {
			return nil, readTokenErr
		}
		token = viper.GetString(tokenKey)
	}

	if clientId == "" {
		viper.SetDefault(clientIdKey, shortuuid.New())
		writeClientIdErr := viper.WriteConfig()
		if writeClientIdErr != nil {
			return nil, writeClientIdErr
		}
		readClientIdErr := viper.ReadInConfig()
		if readClientIdErr != nil {
			return nil, readClientIdErr
		}
		clientId = viper.GetString(clientIdKey)
	}
	return &LocalConfig{Token: token, ClientId: clientId, SchemaVersion: schemaVersion}, nil
}

func (lc *LocalConfigClient) Set(key string, value string) error {
	_, _, _, err := setViperConfig()
	if err != nil {
		return err
	}

	viper.Set(key, value)
	writeClientIdErr := viper.WriteConfig()
	if writeClientIdErr != nil {
		return writeClientIdErr
	}
	return nil
}

func InitLocalConfigFile() error {
	configHome, configName, configType, err := setViperConfig()
	if err != nil {
		return err
	}
	// workaround for creating config file when not exist
	// open issue in viper: https://github.com/spf13/viper/issues/430
	// should be fixed in pr https://github.com/spf13/viper/pull/936
	configPath := filepath.Join(configHome, configName+"."+configType)

	_, err = os.Stat(configHome)
	if err != nil {
		osMkdirErr := os.Mkdir(configHome, os.ModePerm)
		if osMkdirErr != nil {
			return osMkdirErr
		}
	}

	_, err = os.Stat(configPath)
	if err != nil {
		_, osCreateErr := os.Create(configPath)
		if osCreateErr != nil {
			return osCreateErr
		}
	}

	readLocalFileErr := viper.ReadInConfig()
	if readLocalFileErr != nil {
		return readLocalFileErr
	}
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

package localConfig

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/datreeio/datree/pkg/networkValidator"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/lithammer/shortuuid"
	"github.com/spf13/viper"
)

type LocalConfig struct {
	Token         string
	ClientId      string
	SchemaVersion string
	Offline       string
}

type TokenClient interface {
	CreateToken() (*cliClient.CreateTokenResponse, error)
}

type LocalConfigClient struct {
	tokenClient      TokenClient
	networkValidator *networkValidator.NetworkValidator
}

func NewLocalConfigClient(t TokenClient, nv *networkValidator.NetworkValidator) *LocalConfigClient {
	return &LocalConfigClient{
		tokenClient:      t,
		networkValidator: nv,
	}
}

const (
	clientIdKey      = "client_id"
	tokenKey         = "token"
	schemaVersionKey = "schema_version"
	offlineKey       = "offline"
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
	offline := viper.GetString(offlineKey)

	if offline == "" {
		viper.SetDefault(offlineKey, "fail")
		writeOfflineErr := viper.WriteConfig()
		if writeOfflineErr != nil {
			return &LocalConfig{}, writeOfflineErr
		}
		readOfflineErr := viper.ReadInConfig()
		if readOfflineErr != nil {
			return &LocalConfig{}, readOfflineErr
		}
		offline = viper.GetString(offlineKey)
	}

	lc.networkValidator.SetOfflineMode(offline)

	if token == "" {
		createTokenResponse, err := lc.tokenClient.CreateToken()
		if err != nil {
			return &LocalConfig{}, err
		}
		token = createTokenResponse.Token
		viper.SetDefault(tokenKey, token)
		writeTokenErr := viper.WriteConfig()
		if writeTokenErr != nil {
			return &LocalConfig{}, writeTokenErr
		}
		readTokenErr := viper.ReadInConfig()
		if readTokenErr != nil {
			return &LocalConfig{}, readTokenErr
		}
		token = viper.GetString(tokenKey)
	}

	if clientId == "" {
		viper.SetDefault(clientIdKey, shortuuid.New())
		writeClientIdErr := viper.WriteConfig()
		if writeClientIdErr != nil {
			return &LocalConfig{}, writeClientIdErr
		}
		readClientIdErr := viper.ReadInConfig()
		if readClientIdErr != nil {
			return &LocalConfig{}, readClientIdErr
		}
		clientId = viper.GetString(clientIdKey)
	}

	return &LocalConfig{Token: token, ClientId: clientId, SchemaVersion: schemaVersion, Offline: offline}, nil
}

func (lc *LocalConfigClient) Set(key string, value string) error {
	initConfigFileErr := InitLocalConfigFile()
	if initConfigFileErr != nil {
		return initConfigFileErr
	}

	err := validateKeyValueConfig(key, value)
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

	isDirExists, _ := exists(configHome)
	if !isDirExists {
		osMkdirErr := os.Mkdir(configHome, os.ModePerm)
		if osMkdirErr != nil {
			return osMkdirErr
		}
	}

	isConfigExists, _ := exists(configPath)
	if !isConfigExists {
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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

func validateKeyValueConfig(key string, value string) error {
	if key == "offline" && value != "fail" && value != "local" {
		return fmt.Errorf("Invalid offline configuration value- %q\n"+
			"Valid offline values are - fail, local\n", value)
	}
	return nil
}

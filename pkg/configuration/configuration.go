package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Kind string

const (
	KindConfig Kind = "Config"
)

func init() {
	logger.Initialize()
}

type AzureConfig struct {
	Credentials az.Credentials `yaml:"credentials"`
}

type Config struct {
	Version            string      `yaml:"version"`
	Kind               Kind        `yaml:"kind"`
	CurrentEnvironment uuid.UUID   `yaml:"current-environment"`
	AzureConfig        AzureConfig `yaml:"azure-config,omitempty"`
}

//CreateNewEnvironment in Config
func (c *Config) CreateNewEnvironment(name string) (uuid.UUID, error) {
	logger.Debug().Msgf("will try to create environment %s", name)
	env, err := environment.Create(name)
	if err != nil {
		logger.Panic().Err(err).Msg("creation of new environment failed")
	}
	util.EnsureDirectory(path.Join(
		util.UsedEnvironmentDirectory,
		env.Uuid.String(),
		"/shared", //TODO to consts
	))
	c.CurrentEnvironment = env.Uuid
	logger.Debug().Msgf("will try to save updated config %+v", c)
	return env.Uuid, c.Save()
}

//SetUsedEnvironment to another value
func (c *Config) SetUsedEnvironment(u uuid.UUID) error {
	// Check if passed environment id is valid
	isEnvValid, err := environment.IsExisting(u) // TODO think if it should be here
	if err != nil {
		return err
	} else if !isEnvValid {
		return fmt.Errorf("environment %s not found", u.String())
	}

	logger.Debug().Msgf("changing used environment to %s", u.String())
	c.CurrentEnvironment = u
	logger.Debug().Msgf("will try to save updated config %+v", c)
	return c.Save()
}

//GetConfigFilePath from usedConfigFile variable or fails if not set
func (c *Config) GetConfigFilePath() string {
	if util.UsedConfigFile == "" {
		logger.Panic().Msg("variable usedConfigFile not initialized")
	}
	return util.UsedConfigFile
}

//Save Config to usedConfigFile
func (c *Config) Save() error {
	logger.Debug().Msgf("will try to marshal config %+v", c)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("will try to write marshaled data to file %s", util.UsedConfigFile)
	err = ioutil.WriteFile(util.UsedConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) AddAzureCredentials(credentials az.Credentials) {
	c.AzureConfig.Credentials = credentials
}

//GetConfig sets usedConfigFile and usedConfigurationDirectory to default values and returns (existing or just initialized) Config
func GetConfig() (*Config, error) {
	// TODO use SetConfig inside here
	logger.Debug().Msg("will try to get config file")
	if util.UsedConfigurationDirectory == "" {
		util.UsedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	}
	util.EnsureDirectory(util.UsedConfigurationDirectory)

	if util.UsedConfigFile == "" {
		util.UsedConfigFile = path.Join(util.UsedConfigurationDirectory, util.DefaultConfigFileName)
	}

	if util.UsedEnvironmentDirectory == "" {
		util.UsedEnvironmentDirectory = path.Join(util.UsedConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)
	}
	util.EnsureDirectory(util.UsedEnvironmentDirectory)

	if util.UsedRepositoryFile == "" {
		util.UsedRepositoryFile = path.Join(util.UsedConfigurationDirectory, util.DefaultV1RepositoryFileName)
	}

	if util.UsedTempDirectory == "" {
		util.UsedTempDirectory = path.Join(util.UsedConfigurationDirectory, util.DefaultEnvironmentsTempSubdirectory)
	}
	util.EnsureDirectory(util.UsedTempDirectory)

	logger.Debug().Msg("will try to make or get configuration")
	return makeOrGetConfig()
}

//SetConfigDirectory sets variable usedConfigurationDirectory and returns (existing or just initialized) Config
func SetConfigDirectory(configDir string) (*Config, error) {
	return setUsedConfigPaths(configDir, path.Join(configDir, util.DefaultConfigFileName))
}

//setUsedConfigPaths to provided values
func setUsedConfigPaths(configDir string, configFile string) (*Config, error) {
	logger.Debug().Msgf("will try to set config directory to %s", configDir)
	if util.UsedConfigurationDirectory != "" {
		return nil, fmt.Errorf("util.UsedConfigurationDirectory is %s but should be empty on set", util.UsedConfigurationDirectory)
	}
	util.UsedConfigurationDirectory = configDir
	util.EnsureDirectory(util.UsedConfigurationDirectory)

	logger.Debug().Msg("will try to set used config file")
	if util.UsedConfigFile != "" {
		return nil, fmt.Errorf("util.UsedConfigFile is %s but should be empty on set", util.UsedConfigFile)
	}
	util.UsedConfigFile = configFile

	logger.Debug().Msg("will try to set used environments directory")
	if util.UsedEnvironmentDirectory != "" {
		return nil, fmt.Errorf("util.UsedEnvironmentDirectory is %s but should be empty on set", util.UsedEnvironmentDirectory)
	}
	util.UsedEnvironmentDirectory = path.Join(configDir, util.DefaultEnvironmentsSubdirectory)
	util.EnsureDirectory(util.UsedEnvironmentDirectory)

	logger.Debug().Msg("will try to set used temporary directory")
	if util.UsedTempDirectory != "" {
		return nil, fmt.Errorf("util.UsedTempDirectory is %s but should be empty on set", util.UsedTempDirectory)
	}
	util.UsedTempDirectory = path.Join(configDir, util.DefaultEnvironmentsTempSubdirectory)
	util.EnsureDirectory(util.UsedTempDirectory)

	logger.Debug().Msg("will try to set repo config file path")
	if util.UsedRepositoryFile != "" {
		return nil, fmt.Errorf("util.UsedRepositoryFile is %s but should be empty on set", util.UsedRepositoryFile)
	}
	util.UsedRepositoryFile = path.Join(configDir, util.DefaultV1RepositoryFileName)

	logger.Debug().Msg("will try to make or get configuration")
	return makeOrGetConfig()
}

//makeOrGetConfig initializes new config file or reads existing one and returns Config
func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(util.UsedConfigFile); os.IsNotExist(err) {
		logger.Debug().Msg("there is no config file, will try to initialize one")
		config := &Config{
			Version: "v1",
			Kind:    KindConfig,
		}
		err = config.Save()
		if err != nil {
			logger.Panic().Err(err).Msg("failed to save")
		}
		return config, nil
	}
	logger.Debug().Msgf("will try to load existing config file from %s", util.UsedConfigFile)
	config := &Config{}
	logger.Debug().Msgf("trying to open %s file", util.UsedConfigFile)
	file, err := os.Open(util.UsedConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	logger.Debug().Msgf("will try to decode file %s to yaml", util.UsedConfigFile)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

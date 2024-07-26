package configuration

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gmcoringa/tswitch/pkg/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config ... general configuration struct.
type Config struct {
	TerragruntFile string `yaml:"terragruntFile"`
	InstallDir     string `yaml:"installDir"`
	CacheDir       string `yaml:"cacheDir"`
}

var (
	hcl        *string
	configPath *string
	installDir *string
	cacheDir   *string
)

// Load configuration file and parse command line arguments
func Load() (*Config, error) {
	return loadConfig()
}

func InitFlags() {
	var defaultConfigPath = filepath.Join(".tswitch", "config.yaml")

	home, err := os.UserHomeDir()
	if err == nil {
		defaultConfigPath = filepath.Join(home, ".tswitch", "config.yaml")
	}

	configPath = flag.String("config", defaultConfigPath, "absolute path for tswitch configuration")
	hcl = flag.String("terragrunt", "./terragrunt.hcl", "absolute path for root terragrunt.hcl file")
	installDir = flag.String("install_dir", "/usr/local/bin", "directory for binaries installation")
	cacheDir = flag.String("cache_dir", "/usr/local/lib/tswitch", "cache directory for all versions binaries")

	flag.Parse()
}

func loadConfig() (*Config, error) {
	config, err := loadConfigFromDisk()
	if err != nil {
		return nil, err
	}

	// Set defaults for missing configs
	if util.IsBlank(config.InstallDir) {
		config.InstallDir = *installDir
	}

	if util.IsBlank(config.TerragruntFile) {
		config.TerragruntFile = *hcl
	}

	if util.IsBlank(config.CacheDir) {
		config.CacheDir = *cacheDir
	}

	log.Debug("Configuration loaded: ", config)
	return config, nil
}

func loadConfigFromDisk() (*Config, error) {
	log.Info("Loading configuration from file: ", *configPath)

	configFile, err := ioutil.ReadFile(*configPath)
	if os.IsNotExist(err) || err != nil {
		log.Debug("File or not found: ", *configPath)
		log.Debug(err)
		return &Config{}, nil
	} else if err != nil {
		log.Debug("Failed to read configuration from file or not found: ", *configPath)
		return &Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Error("Failed to unmarshal configuration: ", *configPath)
		return &Config{}, err
	}

	return &config, nil
}

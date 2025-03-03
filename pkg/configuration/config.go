package configuration

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gmcoringa/tswitch/pkg/util"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

// Config ... general configuration struct.
type Config struct {
	TerragruntFile string `yaml:"terragruntFile"`
	InstallDir     string `yaml:"installDir"`
	CacheDir       string `yaml:"cacheDir"`
	TerraformImpl  string `yaml:"terraformImplementation"`
}

var (
	hcl        *string
	configPath *string
	installDir *string
	cacheDir   *string
	tfImpl     *string
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
	tfImpl = flag.String("tf_impl", "terraform", "terraform implementation to use, valid values: terraform, tofu")

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

	if util.IsBlank(config.TerraformImpl) {
		config.TerraformImpl = *tfImpl
	}

	if config.TerraformImpl != "terraform" && config.TerraformImpl != "tofu" {
		return nil, fmt.Errorf("invalid terraform implementation: %s. Valid values: terraform | tofu", config.TerraformImpl)
	}

	log.Debug("Configuration loaded: ", config)
	return config, nil
}

func loadConfigFromDisk() (*Config, error) {
	log.Info("Loading configuration from file: ", *configPath)

	configFile, err := os.ReadFile(*configPath)
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

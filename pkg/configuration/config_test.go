package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	hclTest := "/tmp/terragrunt.hcl"
	configPathTest := "./test_config.yaml"
	installDirTest := "/installDir"
	cacheDirTest := "/cacheDir"

	hcl = &hclTest
	configPath = &configPathTest
	installDir = &installDirTest
	cacheDir = &cacheDirTest
}

func TestConfigFileDoesExists(testing *testing.T) {
	validConfig := "./not_found.yaml"
	configPath = &validConfig

	_, err := loadConfig()
	if err != nil {
		assert.Fail(testing, "Failed to parse config", err)
	}
}

func TestParseConfiguration(testing *testing.T) {
	validConfig := "./test_config.yaml"
	configPath = &validConfig

	result, err := loadConfig()

	assert.NoError(testing, err)
	assert.NotNil(testing, result)
}

func TestFailParseWithInvalidConfiguration(testing *testing.T) {
	validConfig := "./test_config_invalid.yaml"
	configPath = &validConfig

	result, err := loadConfig()

	assert.Error(testing, err)
	assert.Nil(testing, result)
}

func TestSetDefaultHclWhenConfigIsMissing(testing *testing.T) {
	validConfig := "./test_config_empty.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, *hcl, config.TerragruntFile)
}

func TestSetHclFromConfig(testing *testing.T) {
	validConfig := "./test_config.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, "./terragrunt.hcl", config.TerragruntFile)
}

func TestSetDefaultInstallDirWhenConfigIsMissing(testing *testing.T) {
	validConfig := "./test_config_empty.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, *installDir, config.InstallDir)
}

func TestSetInstallDirFromConfig(testing *testing.T) {
	validConfig := "./test_config.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, "/tmp/tswitch/bin", config.InstallDir)
}

func TestSetDefaultCacheDirWhenConfigIsMissing(testing *testing.T) {
	validConfig := "./test_config_empty.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, *cacheDir, config.CacheDir)
}

func TestSetCacheDirFromConfig(testing *testing.T) {
	validConfig := "./test_config.yaml"
	configPath = &validConfig
	config, _ := loadConfig()

	assert.Equal(testing, "/tmp/tsiwtch/data", config.CacheDir)
}

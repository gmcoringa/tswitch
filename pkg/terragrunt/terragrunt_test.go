package terragrunt_test

import (
	"os"
	"testing"

	tg "github.com/gmcoringa/tswitch/pkg/terragrunt"
	"github.com/stretchr/testify/assert"
)

func TestListVersions(testing *testing.T) {
	subject := tg.Init()

	result, err := subject.ListVersions()

	assert.NoError(testing, err)
	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, "0.28.0")
}

func TestAddNewVersion(testing *testing.T) {
	err := os.Mkdir("./test_data", 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll("./test_data")

	subject := tg.Init()

	err = subject.AddNewVersion("0.28.0", "./test_data/terragrunt_test_add_new_version")

	assert.NoError(testing, err)
	assert.FileExists(testing, "./test_data/terragrunt_test_add_new_version")
}

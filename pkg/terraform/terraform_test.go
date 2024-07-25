package terraform_test

import (
	"os"
	"testing"

	tf "github.com/gmcoringa/tswitch/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

func TestListVersions(testing *testing.T) {
	subject := tf.Init()

	result, _ := subject.ListVersions()

	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, "1.9.1")
}

func TestAddNewVersion(testing *testing.T) {
	err := os.Mkdir("./test_data", 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll("./test_data")

	subject := tf.Init()

	err = subject.AddNewVersion("1.9.1", "./test_data/terraform_1.9.1")

	assert.NoError(testing, err)
	assert.FileExists(testing, "./test_data/terraform_1.9.1")
}

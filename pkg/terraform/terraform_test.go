package terraform

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTFData = "./test_tf_data"
	tfVersion  = "1.9.1"
)

var tfVersionData = fmt.Sprintf("%s/terraform_%s", testTFData, tfVersion)

func TestTerraformListVersions(testing *testing.T) {
	// Mock getURLContent
	originalGetURLContent := getURLContent
	defer func() { getURLContent = originalGetURLContent }()

	getURLContent = func(url string) ([]string, error) {
		return []string{
			`<a href="/terraform/1.15.0-rc1/">terraform_1.15.0-rc1</a>`,
			`<a href="/terraform/1.15.0-alpha20260218/">terraform_1.15.0-alpha20260218</a>`,
			`<a href="/terraform/1.9.1/">terraform_1.9.1</a>`,
			`<a href="/terraform/1.8.0/">terraform_1.8.0</a>`,
		}, nil
	}

	subject := InitTerraform()

	result, err := subject.ListVersions()

	assert.NoError(testing, err)
	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, tfVersion)
	assert.NotContains(testing, result, "1.15.0")
	assert.NotContains(testing, result, "1.15.0-rc1")
	assert.NotContains(testing, result, "1.15.0-alpha20260218")
}

func TestTerraformAddNewVersion(testing *testing.T) {
	err := os.Mkdir(testTFData, 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll(testTFData)

	subject := InitTerraform()

	err = subject.AddNewVersion(tfVersion, tfVersionData)

	assert.NoError(testing, err)
	assert.FileExists(testing, tfVersionData)
}

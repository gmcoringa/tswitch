package terraform_test

import (
	"fmt"
	"os"
	"testing"

	tf "github.com/gmcoringa/tswitch/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	testTFData = "./test_tf_data"
	tfVersion  = "1.9.1"
)

var tfVersionData = fmt.Sprintf("%s/terraform_%s", testTFData, tfVersion)

func TestTerraformListVersions(testing *testing.T) {
	subject := tf.InitTerraform()

	result, _ := subject.ListVersions()

	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, tfVersion)
}

func TestTerraformAddNewVersion(testing *testing.T) {
	err := os.Mkdir(testTFData, 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll(testTFData)

	subject := tf.InitTerraform()

	err = subject.AddNewVersion(tfVersion, tfVersionData)

	assert.NoError(testing, err)
	assert.FileExists(testing, tfVersionData)
}

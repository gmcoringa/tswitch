package terraform_test

import (
	"fmt"
	"os"
	"testing"

	tf "github.com/gmcoringa/tswitch/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

const (
	testTofuData = "./test_tofu_data"
	tofuVersion  = "1.8.3"
)

var tofuVersionData = fmt.Sprintf("%s/tofu_%s", testTofuData, tofuVersion)

func TestTofuListVersions(testing *testing.T) {
	subject := tf.InitTofu()

	result, _ := subject.ListVersions()

	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, tofuVersion)
}

func TestTofuAddNewVersion(testing *testing.T) {
	err := os.Mkdir(testTofuData, 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll(testTofuData)

	subject := tf.InitTofu()

	err = subject.AddNewVersion(tofuVersion, tofuVersionData)

	assert.NoError(testing, err)
	assert.FileExists(testing, tofuVersionData)
}

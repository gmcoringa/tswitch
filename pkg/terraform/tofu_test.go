package terraform

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

const (
	testTofuData = "./test_tofu_data"
	tofuVersion  = "1.8.3"
)

var tofuVersionData = fmt.Sprintf("%s/tofu_%s", testTofuData, tofuVersion)

func TestTofuListVersions(testing *testing.T) {
	// Mock listRemoteTagsTofu
	originalListRemoteTagsTofu := listRemoteTagsTofu
	defer func() { listRemoteTagsTofu = originalListRemoteTagsTofu }()

	listRemoteTagsTofu = func(repoURL string) ([]*plumbing.Reference, error) {
		return []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/tags/v1.6.0-alpha1", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.7.0-beta1", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.8.0-rc1", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.8.3", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.9.0", "hash"),
		}, nil
	}

	subject := InitTofu()

	result, err := subject.ListVersions()

	assert.NoError(testing, err)
	assert.NotEmpty(testing, result)
	assert.Contains(testing, result, tofuVersion)
	assert.Contains(testing, result, "1.9.0")
	assert.NotContains(testing, result, "1.6.0")
	assert.NotContains(testing, result, "1.6.0-alpha1")
	assert.NotContains(testing, result, "1.7.0-beta1")
	assert.NotContains(testing, result, "1.8.0-rc1")
}

func TestTofuAddNewVersion(testing *testing.T) {
	err := os.Mkdir(testTofuData, 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll(testTofuData)

	subject := InitTofu()

	err = subject.AddNewVersion(tofuVersion, tofuVersionData)

	assert.NoError(testing, err)
	assert.FileExists(testing, tofuVersionData)
}

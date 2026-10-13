package terragrunt

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestListVersions(testing *testing.T) {
	// Mock listRemoteTagsTerragrunt
	originalListRemoteTagsTerragrunt := listRemoteTagsTerragrunt
	defer func() { listRemoteTagsTerragrunt = originalListRemoteTagsTerragrunt }()

	listRemoteTagsTerragrunt = func(repoURL string) ([]*plumbing.Reference, error) {
		return []*plumbing.Reference{
			plumbing.NewReferenceFromStrings("refs/tags/alpha-2026021301", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/beta-2025111005", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/alpha2025022701", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v1.0.0-rc1", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.64.0", "hash"),
			plumbing.NewReferenceFromStrings("refs/tags/v0.99.5", "hash"),
			plumbing.NewReferenceFromStrings("refs/heads/main", "hash"),
		}, nil
	}

	subject := Init()

	result, err := subject.ListVersions()

	assert.NoError(testing, err)
	assert.ElementsMatch(testing, []string{"0.64.0", "0.99.5"}, result)
}

func TestAddNewVersion(testing *testing.T) {
	err := os.Mkdir("./test_data", 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll("./test_data")

	subject := Init()

	err = subject.AddNewVersion("0.64.0", "./test_data/terragrunt_test_add_new_version")

	assert.NoError(testing, err)
	assert.FileExists(testing, "./test_data/terragrunt_test_add_new_version")
}

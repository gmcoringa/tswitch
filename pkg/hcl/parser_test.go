package hcl_test

import (
	"testing"

	"github.com/gmcoringa/tswitch/pkg/hcl"
	"github.com/stretchr/testify/assert"
)

const (
	validFile   string = "terragrunt.hcl"
	invalidFile string = "terragrunt-no-constraints.hcl"
)

func TestParseConfiguration(testing *testing.T) {
	_, err := hcl.Parse(validFile)

	if err != nil {
		testing.Fatalf(`Parse(file) failed with error %v`, err)
	}
}

func TestReadTerraformVersion(testing *testing.T) {
	result, _ := hcl.Parse(validFile)

	assert.Equal(testing, ">= 0.13, < 0.14", result.Terraform)
}

func TestReadTerragruntVersion(testing *testing.T) {
	result, _ := hcl.Parse(validFile)

	assert.Equal(testing, ">= 0.26, < 0.27", result.Terragrunt)
}

func TestFailWithMissingConstraints(testing *testing.T) {
	_, err := hcl.Parse(invalidFile)

	if err == nil {
		testing.Fatalf(`Should fail with missing constraints in file %v`, invalidFile)
	}
}

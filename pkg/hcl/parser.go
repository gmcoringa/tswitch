package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	log "github.com/sirupsen/logrus"
)

type VersionConstraints struct {
	Terraform  string   `hcl:"terraform_version_constraint"`
	Terragrunt string   `hcl:"terragrunt_version_constraint"`
	Remain     hcl.Body `hcl:",remain"`
}

// Parse the fiven file to retrieve terraform and terragrunt constraints
func Parse(filename string) (*VersionConstraints, error) {
	var versionConstraints VersionConstraints
	err := hclsimple.DecodeFile(filename, nil, &versionConstraints)
	if err != nil {
		log.Error("Failed to decode hcl file ", filename)
		return nil, err
	}

	return &versionConstraints, nil
}

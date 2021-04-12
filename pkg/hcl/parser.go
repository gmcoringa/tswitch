package hcl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gmcoringa/tswitch/pkg/util"
	log "github.com/sirupsen/logrus"
)

type VersionConstraints struct {
	Terraform  string
	Terragrunt string
}

// Parse the fiven file to retrieve terraform and terragrunt constraints
func Parse(filename string) (*VersionConstraints, error) {
	config := &VersionConstraints{}

	if len(filename) == 0 {
		return config, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Error("Failed to load hcl file", filename)
		return config, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "terragrunt_version_constraint") {
			config.Terragrunt = retrieveValue(line)
		}

		if strings.HasPrefix(line, "terraform_version_constraint") {
			config.Terraform = retrieveValue(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Failed to parse hcl file", filename)
		return config, err
	}

	if util.IsBlank(config.Terraform) || util.IsBlank(config.Terragrunt) {
		return config, fmt.Errorf("constraints for terraform or terragrunt missing in %s", filename)
	}

	return config, nil
}

// Retrieve the value of line with the format key=value
func retrieveValue(line string) string {
	if equal := strings.Index(line, "="); equal >= 0 {
		if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
			value := ""

			if len(line) > equal {
				value = strings.TrimSpace(line[equal+1:])
			}

			return strings.ReplaceAll(value, "\"", "")
		}
	}

	return ""
}

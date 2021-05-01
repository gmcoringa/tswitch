package terragrunt

import (
	"fmt"
	"runtime"
	"strings"

	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/gmcoringa/tswitch/pkg/lib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

type Terragrunt struct {
	repoURL        string
	downloadURL    string
	name           string
	installVersion string
}

func Init() lib.Resolver {
	return Terragrunt{
		repoURL:        "https://github.com/gruntwork-io/terragrunt",
		downloadURL:    "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%s_%s",
		name:           "terragrunt",
		installVersion: "terragrunt_",
	}
}

func (tg Terragrunt) Name() string {
	return tg.name
}

// ListVersions :  Get the list of available terragrunt versions
func (tg Terragrunt) ListVersions() ([]string, error) {
	// Create the remote with repository URL
	repo := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{tg.repoURL},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := repo.List(&git.ListOptions{})
	if err != nil {
		log.Error("Failed to retrieve version list for terragrunt")
		return nil, err
	}

	// Filters the references list and only keeps tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, strings.TrimPrefix(ref.Name().Short(), "v"))
		}
	}

	return tags, nil
}

// AddNewVersion: download the given version of terragrunt and copy the executable binary
// to the given destination
func (tg Terragrunt) AddNewVersion(version string, destination string) error {
	url := fmt.Sprintf(tg.downloadURL, version, runtime.GOOS, runtime.GOARCH)
	err := lio.DownloadFromURLToLocation(destination, url)
	if err != nil {
		log.Error("Failed to download terragrunt binary from ", url, " to ", destination)
	}

	return err
}

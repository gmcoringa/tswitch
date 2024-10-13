package terraform

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/gmcoringa/tswitch/pkg/lib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

type Tofu struct {
	repoURL        string
	downloadURL    string
	name           string
	implementation string
}

func InitTofu() lib.Resolver {
	return Tofu{
		repoURL:        "https://github.com/opentofu/opentofu",
		downloadURL:    "https://github.com/opentofu/opentofu/releases/download/v%s/tofu_%s_%s_%s.tar.gz",
		name:           "terraform",
		implementation: "tofu",
	}
}

func (tf Tofu) Name() string {
	return tf.name
}

func (tf Tofu) Implementation() string {
	return tf.implementation
}

// ListVersions :  Get the list of available tofu versions
func (tf Tofu) ListVersions() ([]string, error) {
	// Create the remote with repository URL
	repo := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{tf.repoURL},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := repo.List(&git.ListOptions{})
	if err != nil {
		log.Error("Failed to retrieve version list for tofu")
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

// AddNewVersion: download the given version of tofu and copy the executable binary
// to the given destination
func (tf Tofu) AddNewVersion(version string, destination string) error {
	url := fmt.Sprintf(tf.downloadURL, version, version, runtime.GOOS, runtime.GOARCH)

	tempDir, err := os.MkdirTemp(os.TempDir(), "tswitch-")
	if err != nil {
		log.Error("Failed to create temporaly directory for download binaries in ", tempDir)
		return err
	}

	defer os.RemoveAll(tempDir)
	tarFile := filepath.Join(tempDir, "tofu.tar.gz")

	err = lio.DownloadFromURLToLocation(tarFile, url)
	if err != nil {
		return err
	}

	// Extract the tar file
	tarFileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer tarFileReader.Close()

	uncompressLocation := filepath.Join(tempDir, "tofu")
	err = extractTarGz(tarFileReader, uncompressLocation)
	if err != nil {
		return err
	}

	source := filepath.Join(tempDir, "tofu", "tofu")
	err = lio.Move(source, destination)
	if err != nil {
		log.Error("Failed to move ", source, " -> ", destination)
		return err
	}

	return err
}

// extractTarGz extracts a .tar.gz file to the specified destination
func extractTarGz(gzipStream io.Reader, destination string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)
	// Create destination directory
	if err := os.MkdirAll(destination, os.FileMode(os.ModePerm)); err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create file
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return err
			}
			file.Close()
		default:
			// Unsupported type
			return fmt.Errorf("unsupported type: %v in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}

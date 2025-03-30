package terraform

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/gmcoringa/tswitch/pkg/lib"
	log "github.com/sirupsen/logrus"
)

type Terraform struct {
	hashiURL       string
	name           string
	installVersion string
}

func InitTerraform() lib.Resolver {
	return Terraform{
		hashiURL:       "https://releases.hashicorp.com/terraform/",
		name:           "terraform",
		installVersion: "terraform_",
	}
}

func (tf Terraform) Name() string {
	return tf.name
}

func (tf Terraform) Implementation() string {
	return tf.name
}

// ListVersions :  Get the list of available terraform versions
func (tf Terraform) ListVersions() ([]string, error) {
	result, err := getURLContent(tf.hashiURL)
	if err != nil {
		return nil, err
	}

	var versions []string
	regex := regexp.MustCompile(`\/(\d+\.\d+\.\d+)\/?`)

	for item := range result {
		if regex.MatchString(result[item]) {
			str := regex.FindString(result[item])
			trimstr := strings.Trim(str, "/\"")
			versions = append(versions, trimstr)
		}
	}

	return versions, nil
}

// AddNewVersion: download the given version of terraform and copy the executable binary
// to the given destination
func (tf Terraform) AddNewVersion(version string, destination string) error {
	url := tf.hashiURL + version + "/" + tf.installVersion + version + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".zip"

	tempDir, err := os.MkdirTemp(os.TempDir(), "tswitch-")
	if err != nil {
		log.Error("Failed to create temporaly directory for download binaries in ", tempDir)
		return err
	}

	defer os.RemoveAll(tempDir)
	zipFile := filepath.Join(tempDir, "terraform.zip")

	err = lio.DownloadFromURLToLocation(zipFile, url)
	if err != nil {
		return err
	}

	/* unzip the downloaded zipfile */
	uncompressLocation := filepath.Join(tempDir, "terraform")
	err = unzip(zipFile, uncompressLocation)
	if err != nil {
		return err
	}

	source := filepath.Join(tempDir, "terraform", "terraform")
	err = lio.Move(source, destination)
	if err != nil {
		log.Error("Failed to move ", source, " -> ", destination)
		return err
	}

	return nil
}

// getURLContent : get the content of the given url
func getURLContent(url string) ([]string, error) { //nolint: gosec
	response, err := http.Get(url) //nolint: gosec
	if err != nil {
		log.Error("Failed to access url ", url)
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Error reading body")
		return nil, err
	}

	bodyString := string(body)
	result := strings.Split(bodyString, "\n")

	return result, nil
}

// Decompress a zip archive, moving all files and folders to the give destination
func unzip(zipFile string, destination string) error {
	log.Debug("Decompressing ", zipFile, "on destination ", destination)
	reader, err := zip.OpenReader(zipFile)

	if err != nil {
		log.Error("Failed to zip file ", zipFile)
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		log.Debug("Moving files from zip to destination ", destination)

		fileReader, err := file.Open()
		if err != nil {
			log.Error("Failed to read file ", file)
			return err
		}

		defer fileReader.Close()
		fpath := filepath.Join(destination, file.Name) //nolint: gosec
		fpath = filepath.Clean(fpath)

		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			log.Error("Invalid file path: ", fpath)
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				log.Error("Failed to create directory ", fpath)
				return err
			}
		} else {
			dir := filepath.Dir(fpath)
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Error("Failed to create directory ", dir)
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				log.Error("Failed to read file ", fpath)
				return err
			}

			_, err = io.Copy(outFile, fileReader) //nolint: gosec
			outFile.Close()
			if err != nil {
				log.Error("Failed to copy file ", outFile)
				return err
			}
		}
	}

	return nil
}

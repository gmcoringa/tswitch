package io

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// downloadFromURLToLocation : download data from the given url to the given file location
func DownloadFromURLToLocation(location string, url string) error { //nolint: gosec
	log.Info("Downloading ", url, " to ", location)

	response, err := http.Get(url) //nolint: gosec
	if err != nil {
		log.Error("Failed to download from ", url)
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("failed to download version from %s, please check if content exists", url)
	}

	output, err := os.Create(location)
	if err != nil {
		log.Error("Failed to create file ", location)
		return err
	}

	defer output.Close()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Error("Failed to save response from ", url, " into ", location)
		return err
	}

	return nil
}

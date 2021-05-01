package io_test

import (
	"os"
	"testing"

	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFromURLToLocation(testing *testing.T) {
	err := lio.CreateDirIfNotExist("./test_data")
	assert.NoError(testing, err, "Requirement failed")
	defer os.Remove("./test_data/test_download")

	err = lio.DownloadFromURLToLocation("./test_data/test_download", "https://github.com")

	assert.NoError(testing, err)
	assert.FileExists(testing, "./test_data/test_download")
}

func TestDownloadFromURLToLocationShouldFailWithNon200Code(testing *testing.T) {
	err := lio.DownloadFromURLToLocation("./test_data/test_download_non_200", "https://github.com/gmcoringa/tswitch/blob/main/NOT_FOUND")

	assert.Error(testing, err)
}

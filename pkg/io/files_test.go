package io_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/stretchr/testify/assert"
)

func TestForceSymLink(testing *testing.T) {
	err := lio.CreateDirIfNotExist("./test_data")
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll("./test_data")

	link := "./test_data/force_link"
	_, filename, _, _ := runtime.Caller(0)
	_ = os.Remove(link) // be sure the link do not exists

	lio.SymLink(link, filename)
	lio.ForceSymLink(link, filename)

	assert.True(testing, lio.IsSymlink(link))
}

func TestCreateDirIfNotExistWhenDirExists(testing *testing.T) {
	path := "./test_data/create_dir_exists"
	err := os.MkdirAll(path, 0755)
	assert.NoError(testing, err, "Requirement failed")
	defer os.RemoveAll(path)

	err = lio.CreateDirIfNotExist(path)

	assert.NoError(testing, err)
}

func TestCreateDirIfNotExistWhenDirDoesNotExists(testing *testing.T) {
	path := "./test_data/create_dir_does_not_exists"

	err := lio.CreateDirIfNotExist(path)
	defer os.RemoveAll(path)

	assert.NoError(testing, err)
	assert.DirExists(testing, path)
}

func TestRemoveFileIfExistWhenFileDoesNotExists(testing *testing.T) {
	path := "./test_data/remove_file_does_not_exists"

	err := lio.RemoveFileIfExist(path)

	assert.NoError(testing, err)
	assert.NoFileExists(testing, path)
}

func TestRemoveFileIfExistWhenFileExists(testing *testing.T) {
	path := "./test_data/remove_file_exists"
	_, err := os.Create(path)
	assert.NoError(testing, err, "Requirement failed")

	err = lio.RemoveFileIfExist(path)

	assert.NoError(testing, err)
	assert.NoFileExists(testing, path)
}

func TestRemoveDirIfExistWhenDirDoesNotExists(testing *testing.T) {
	path := "./test_data/remove_dir_does_not_exists"

	err := lio.RemoveDirIfExist(path)

	assert.NoError(testing, err)
	assert.NoDirExists(testing, path)
}

func TestRemoveDirIfExistWhenDirExists(testing *testing.T) {
	path := "./test_data/remove_dir_exists"
	err := lio.CreateDirIfNotExist(path)
	assert.NoError(testing, err, "Requirement failed")

	err = lio.RemoveDirIfExist(path)

	assert.NoError(testing, err)
	assert.NoDirExists(testing, path)
}

func TestSetExecutableShouldFaildWithFileNotFound(testing *testing.T) {
	err := lio.SetExecutable("./non_existent_file")

	assert.Error(testing, err)
}

func TestSetExecutableShouldFaildWithDirectories(testing *testing.T) {
	_, filename, _, _ := runtime.Caller(0)

	err := lio.SetExecutable(filepath.Dir(filename))

	assert.Error(testing, err)
}

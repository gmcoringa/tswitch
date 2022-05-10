package lib_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gmcoringa/tswitch/mocks"
	"github.com/gmcoringa/tswitch/pkg/configuration"
	"github.com/gmcoringa/tswitch/pkg/db"
	lio "github.com/gmcoringa/tswitch/pkg/io"
	"github.com/gmcoringa/tswitch/pkg/lib"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	config = &configuration.Config{
		CacheDir:   "test_data",
		InstallDir: "test_data/bin",
	}
)

func TestCurrentVersionInUse(testing *testing.T) {
	controller := gomock.NewController(testing)
	defer controller.Finish()

	path := "test_data/test_current_version/source"
	constraint := ">=1.0, <1.1"
	target := "test_current_version"
	currentVersion := &db.BinVersion{
		Version: "1.0.2",
		Path:    path,
	}

	err := lio.CreateDirIfNotExist(config.InstallDir)
	assert.NoError(testing, err, "Requirement for test failed")
	err = lio.CreateDirIfNotExist(filepath.Dir(path))
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(path)
	assert.NoError(testing, err, "Requirement for test failed")
	defer os.Remove(path)

	mockDb := mocks.NewMockDatabase(controller)
	mockResolver := mocks.NewMockResolver(controller)
	mockResolver.EXPECT().Name().Return(target).AnyTimes()
	mockResolver.EXPECT().ListVersions().Return([]string{"1.1.0", "1.0.2", "1.0.1"}, nil)
	mockDb.EXPECT().GetCurrent(gomock.Eq(target)).Return(currentVersion, nil)

	mockDb.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
	mockDb.EXPECT().SetCurrent(gomock.Any(), gomock.Any()).Times(0)
	mockDb.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockResolver.EXPECT().AddNewVersion(gomock.Any(), gomock.Any()).Times(0)

	subject := lib.CreateInstaller(config, mockDb, mockResolver)

	subject.Install(constraint)
}

// nolint:gocritic
func TestVersionExistsAndIsNotCurrent(testing *testing.T) {
	controller := gomock.NewController(testing)
	defer controller.Finish()

	target := "test_version_exists_not_current"
	path := filepath.Join(config.CacheDir, target, "source")
	constraint := ">=1.0, <1.1"
	version := "1.0.2"
	currentVersion := &db.BinVersion{
		Version: "1.0.1",
		Path:    path,
	}
	latestVersion := &db.BinVersion{
		Version: "1.0.2",
		Path:    path,
	}

	err := lio.CreateDirIfNotExist(filepath.Join(config.InstallDir))
	assert.NoError(testing, err, "Requirement for test failed")
	err = lio.CreateDirIfNotExist(filepath.Join(config.CacheDir, target))
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(path)
	assert.NoError(testing, err, "Requirement for test failed")
	defer os.Remove(filepath.Join(config.InstallDir, target))
	defer os.Remove(filepath.Join(config.InstallDir, target))

	mockDb := mocks.NewMockDatabase(controller)
	mockResolver := mocks.NewMockResolver(controller)
	mockResolver.EXPECT().Name().Return(target).AnyTimes()
	mockResolver.EXPECT().ListVersions().Return([]string{"1.1.0", "1.0.2", "1.0.1"}, nil)
	mockDb.EXPECT().GetCurrent(gomock.Eq(target)).Return(currentVersion, nil)
	mockDb.EXPECT().Get(target, version).Return(latestVersion, nil)

	mockDb.EXPECT().SetCurrent(target, version).Times(1)
	mockDb.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockResolver.EXPECT().AddNewVersion(gomock.Any(), gomock.Any()).Times(0)

	subject := lib.CreateInstaller(config, mockDb, mockResolver)

	subject.Install(constraint)
}

// nolint:dupl
func TestVersionDoesNotExists(testing *testing.T) {
	controller := gomock.NewController(testing)
	defer controller.Finish()

	target := "test_version_does_not_exists"
	path := filepath.Join(config.CacheDir, target, target)
	constraint := ">=1.0, <1.1"
	version := "1.0.2"
	currentVersion := &db.BinVersion{
		Version: "1.0.1",
		Path:    path + "_1.0.1",
	}
	newVersion := &db.BinVersion{
		Version: version,
		Path:    path + "_1.0.2",
	}

	err := lio.CreateDirIfNotExist(filepath.Join(config.CacheDir, target))
	assert.NoError(testing, err, "Requirement for test failed")
	err = lio.CreateDirIfNotExist(config.InstallDir)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(currentVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(newVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	defer os.Remove(filepath.Join(config.CacheDir, target))
	defer os.Remove(filepath.Join(config.InstallDir, target))

	mockDb := mocks.NewMockDatabase(controller)
	mockResolver := mocks.NewMockResolver(controller)
	mockResolver.EXPECT().Name().Return(target).AnyTimes()
	mockResolver.EXPECT().ListVersions().Return([]string{"1.1.0", "1.0.2", "1.0.1"}, nil)
	mockDb.EXPECT().GetCurrent(gomock.Eq(target)).Return(currentVersion, nil)
	mockDb.EXPECT().Get(target, version).Return(nil, fmt.Errorf("Not found"))

	mockDb.EXPECT().SetCurrent(gomock.Any(), gomock.Any()).Times(0)
	mockDb.EXPECT().Add(target, newVersion, true).Times(1)
	mockResolver.EXPECT().AddNewVersion(version, newVersion.Path).Times(1)

	subject := lib.CreateInstaller(config, mockDb, mockResolver)

	subject.Install(constraint)
}

// nolint:dupl
func TestVersionMinorVersionConstraint(testing *testing.T) {
	controller := gomock.NewController(testing)
	defer controller.Finish()

	target := "test_version_minor_constraint"
	path := filepath.Join(config.CacheDir, target, target)
	constraint := "~1.0"
	version := "1.0.2"
	currentVersion := &db.BinVersion{
		Version: "1.0.1",
		Path:    path + "_1.0.1",
	}
	newVersion := &db.BinVersion{
		Version: version,
		Path:    path + "_1.0.2",
	}

	err := lio.CreateDirIfNotExist(filepath.Join(config.CacheDir, target))
	assert.NoError(testing, err, "Requirement for test failed")
	err = lio.CreateDirIfNotExist(config.InstallDir)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(currentVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(newVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	defer os.Remove(filepath.Join(config.CacheDir, target))
	defer os.Remove(filepath.Join(config.InstallDir, target))

	mockDb := mocks.NewMockDatabase(controller)
	mockResolver := mocks.NewMockResolver(controller)
	mockResolver.EXPECT().Name().Return(target).AnyTimes()
	mockResolver.EXPECT().ListVersions().Return([]string{"1.1.0", "1.0.2", "1.0.1"}, nil)
	mockDb.EXPECT().GetCurrent(gomock.Eq(target)).Return(currentVersion, nil)
	mockDb.EXPECT().Get(target, version).Return(nil, fmt.Errorf("Not found"))

	mockDb.EXPECT().SetCurrent(gomock.Any(), gomock.Any()).Times(0)
	mockDb.EXPECT().Add(target, newVersion, true).Times(1)
	mockResolver.EXPECT().AddNewVersion(version, newVersion.Path).Times(1)

	subject := lib.CreateInstaller(config, mockDb, mockResolver)

	subject.Install(constraint)
}

// nolint:dupl
func TestVersionMajorVersionConstraint(testing *testing.T) {
	controller := gomock.NewController(testing)
	defer controller.Finish()

	target := "test_version_major_constraint"
	path := filepath.Join(config.CacheDir, target, target)
	constraint := "~1"
	version := "1.1.0"
	currentVersion := &db.BinVersion{
		Version: "1.0.1",
		Path:    path + "_1.0.1",
	}
	newVersion := &db.BinVersion{
		Version: version,
		Path:    path + "_1.1.0",
	}

	err := lio.CreateDirIfNotExist(filepath.Join(config.CacheDir, target))
	assert.NoError(testing, err, "Requirement for test failed")
	err = lio.CreateDirIfNotExist(config.InstallDir)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(currentVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	_, err = os.Create(newVersion.Path)
	assert.NoError(testing, err, "Requirement for test failed")
	defer os.Remove(filepath.Join(config.CacheDir, target))
	defer os.Remove(filepath.Join(config.InstallDir, target))

	mockDb := mocks.NewMockDatabase(controller)
	mockResolver := mocks.NewMockResolver(controller)
	mockResolver.EXPECT().Name().Return(target).AnyTimes()
	mockResolver.EXPECT().ListVersions().Return([]string{"1.1.0", "1.0.2", "1.0.1"}, nil)
	mockDb.EXPECT().GetCurrent(gomock.Eq(target)).Return(currentVersion, nil)
	mockDb.EXPECT().Get(target, version).Return(nil, fmt.Errorf("Not found"))

	mockDb.EXPECT().SetCurrent(gomock.Any(), gomock.Any()).Times(0)
	mockDb.EXPECT().Add(target, newVersion, true).Times(1)
	mockResolver.EXPECT().AddNewVersion(version, newVersion.Path).Times(1)

	subject := lib.CreateInstaller(config, mockDb, mockResolver)

	subject.Install(constraint)
}

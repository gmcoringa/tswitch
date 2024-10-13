package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	semver "github.com/Masterminds/semver/v3"
	"github.com/gmcoringa/tswitch/pkg/configuration"
	"github.com/gmcoringa/tswitch/pkg/db"
	lio "github.com/gmcoringa/tswitch/pkg/io"
	log "github.com/sirupsen/logrus"
)

type Installer interface {
	Install(constraint string)
}

type Resolver interface {
	ListVersions() ([]string, error)
	Name() string
	Implementation() string
	AddNewVersion(version string, destination string) error
}

type LocalInstaller struct {
	Config  *configuration.Config
	DB      db.Database
	BinPath string
	Target  Resolver
}

func CreateInstaller(config *configuration.Config, database db.Database, installer Resolver) Installer {
	return &LocalInstaller{
		Config:  config,
		DB:      database,
		BinPath: filepath.Join(config.InstallDir, installer.Name()),
		Target:  installer,
	}
}

// Install : Install the provided version in the argument
func (installer *LocalInstaller) Install(constraint string) {
	err := lio.CreateDirIfNotExist(filepath.Join(installer.Config.CacheDir, installer.Target.Implementation()))
	if err != nil {
		log.Error("Failed to create cache dir for binaries: ", filepath.Join(installer.Config.CacheDir, installer.Target.Implementation()))
		log.Fatal(err)
	}

	err = lio.CreateDirIfNotExist(installer.Config.InstallDir)
	if err != nil {
		log.Error("Failed to create directory for binaries usage: ", installer.Config.InstallDir)
		log.Fatal(err)
	}

	versionList, err := installer.Target.ListVersions()
	if err != nil {
		log.Errorln("Error looking for available versions")
		log.Errorln(err)
		os.Exit(1)
	}

	version, err := findVersion(constraint, versionList)
	if err != nil {
		log.Error("Failed to find matching version for ", installer.Target.Implementation(), " with constraint ", constraint)
		log.Error(err)
		os.Exit(1)
	}

	// Current version already in use
	currentInstall, err := installer.DB.GetCurrent(installer.Target.Implementation())
	if err == nil && currentInstall.Version == version {
		// Force symbolic link for requested version
		lio.ForceSymLink(installer.BinPath, currentInstall.Path)
		logInstallation(installer.Target.Implementation(), version)
		return
	}

	// Version already installed, just set as current
	versionInstall, err := installer.DB.Get(installer.Target.Implementation(), version)
	if err == nil {
		log.Debug(installer.Target.Implementation(), "version ", version, "found, setting as current version")
		err = installer.DB.SetCurrent(installer.Target.Implementation(), version)
		if err != nil {
			log.Warn("Failed to update db [", installer.Target.Implementation(), "] with current version ", version, ", ignoring")
		}

		lio.ForceSymLink(installer.BinPath, versionInstall.Path)
		logInstallation(installer.Target.Implementation(), version)
		return
	}

	// Version not cached, download needed
	binName := fmt.Sprintf("%s_%s", installer.Target.Implementation(), version)
	destination := filepath.Join(installer.Config.CacheDir, installer.Target.Implementation(), binName)
	err = installer.Target.AddNewVersion(version, destination)
	if err != nil {
		log.Error("Failed to set new version [", version, "] on destination ", destination)
		log.Error(err)
		os.Exit(1)
	}
	installment := db.BinVersion{
		Version: version,
		Path:    destination,
	}

	err = installer.DB.Add(installer.Target.Implementation(), &installment, true)
	if err != nil {
		log.Error("Failed to set new version for ", installer.Target.Implementation(), " in DB, version info: ", installment)
		log.Error(err)
		os.Exit(1)
	}

	lio.ForceSymLink(installer.BinPath, installment.Path)
	err = lio.SetExecutable(installment.Path)
	if err != nil {
		log.Error("Failed to change permissions of ", installment.Path)
		log.Error(err)
		os.Exit(1)
	}

	logInstallation(installer.Target.Implementation(), version)
}

func logInstallation(target string, version string) {
	log.Info("Installed ", target, " version: ", version)
}

func findVersion(constraint string, versionList []string) (string, error) {
	constraints, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Error("Invalid constraint: ", constraint)
		return "", err
	}

	versions := make([]*semver.Version, len(versionList))
	for index, item := range versionList {
		version, err := semver.NewVersion(item)

		if err != nil {
			log.Error("Failed to parse version ", version)
			return "", err
		}

		versions[index] = version
	}

	sort.Sort(sort.Reverse(semver.Collection(versions)))

	for _, candidate := range versions {
		if constraints.Check(candidate) { // Validate a version against a constraint
			version := candidate.String()
			log.Info("Found match version [", version, "] for constraint [", constraint, "]")

			return version, nil
		}
	}

	return "", fmt.Errorf("no version match for constraint %s", constraint)
}

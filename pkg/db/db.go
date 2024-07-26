package db

import (
	"path/filepath"

	"github.com/gmcoringa/tswitch/pkg/configuration"
	scribble "github.com/nanobox-io/golang-scribble"
	log "github.com/sirupsen/logrus"
)

const (
	currentVersion = "CURRENT_VERSION"
)

type BinVersion struct {
	Version string
	Path    string
}

type Database interface {
	Get(target string, version string) (*BinVersion, error)
	GetCurrent(target string) (*BinVersion, error)
	SetCurrent(target string, version string) error
	Add(target string, installment *BinVersion, current bool) error
}

type LocalDatabase struct {
	driver *scribble.Driver
}

func Init(config *configuration.Config) (Database, error) {
	filepath.Join(config.CacheDir, "db")
	db, err := scribble.New(filepath.Join(config.CacheDir, "db"), nil)
	localDB := LocalDatabase{driver: db}

	return localDB, err
}

func (db LocalDatabase) Add(target string, installment *BinVersion, current bool) error {
	err := db.driver.Write(target, installment.Version, installment)

	if err != nil {
		log.Error("Failed to save installment ", *installment)
		return err
	}

	if current {
		err := db.driver.Write(target, currentVersion, installment)

		if err != nil {
			log.Error("Failed to set current version for ", *installment)
		}
	}

	return err
}

func (db LocalDatabase) Get(target string, version string) (*BinVersion, error) {
	binVersion := &BinVersion{}
	err := db.driver.Read(target, version, binVersion)

	return binVersion, err
}

func (db LocalDatabase) GetCurrent(target string) (*BinVersion, error) {
	binVersion := &BinVersion{}
	err := db.driver.Read(target, currentVersion, binVersion)

	return binVersion, err
}

func (db LocalDatabase) SetCurrent(target string, version string) error {
	binVersion := &BinVersion{}
	err := db.driver.Read(target, version, binVersion)
	if err != nil {
		log.Error("Cannot find installment for target [", target, "] and version [", version, "]")
		return err
	}

	return db.driver.Write(target, currentVersion, binVersion)
}

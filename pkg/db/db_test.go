package db_test

import (
	"testing"

	"github.com/gmcoringa/tswitch/pkg/configuration"
	"github.com/gmcoringa/tswitch/pkg/db"
	"github.com/stretchr/testify/assert"
)

var database db.Database

func init() {
	config := &configuration.Config{
		CacheDir: "./test_data",
	}

	database, _ = db.Init(config)
}

func TestAddWithCurrent(testing *testing.T) {
	target := "testAdd"
	add := db.BinVersion{
		Version: "1.0.0",
		Path:    "/v1",
	}

	err := database.Add(target, &add, true)
	assert.NoError(testing, err)

	added, _ := database.Get(target, "1.0.0")
	current, _ := database.GetCurrent(target)

	assert.Equal(testing, &add, added)
	assert.Equal(testing, &add, current)
}

func TestAddWithoutCurrent(testing *testing.T) {
	target := "testAddNoCurrent"
	add := db.BinVersion{
		Version: "1.0.0",
		Path:    "/v1",
	}

	err := database.Add(target, &add, false)
	assert.NoError(testing, err)

	added, _ := database.Get(target, "1.0.0")
	_, err = database.GetCurrent(target)

	assert.Equal(testing, &add, added)
	assert.NotNil(testing, err, "Should not exists current version")
}

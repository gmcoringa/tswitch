package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gmcoringa/tswitch/pkg/configuration"
	"github.com/gmcoringa/tswitch/pkg/db"
	"github.com/gmcoringa/tswitch/pkg/hcl"
	"github.com/gmcoringa/tswitch/pkg/lib"
	formatter "github.com/gmcoringa/tswitch/pkg/log"
	tf "github.com/gmcoringa/tswitch/pkg/terraform"
	tg "github.com/gmcoringa/tswitch/pkg/terragrunt"
	log "github.com/sirupsen/logrus"
)

var (
	version = "snapshot"
)

func main() {

	formatter := formatter.Format{}
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stderr)
	logLevel := flag.String("log_level", "info", "logging level, valid vaulues: trace, debug, info, warn, error, fatal, panic")
	displayVersion := flag.Bool("v", false, "prints current roxy version")
	// init other flags
	configuration.InitFlags()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)

	if len(os.Args) <= 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *displayVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	configuration, err := configuration.Load()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	constraints, err := hcl.Parse(configuration.TerragruntFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	db, err := db.Init(configuration)
	if err != nil {
		log.Error("Failed to load local db", err)
		os.Exit(1)
	}

	tfResolver := tf.Init()
	terraform := lib.CreateInstaller(configuration, db, tfResolver)
	terraform.Install(constraints.Terraform)

	tgResolver := tg.Init()
	terragrunt := lib.CreateInstaller(configuration, db, tgResolver)
	terragrunt.Install(constraints.Terragrunt)
}

package terraform

import (
	"github.com/gmcoringa/tswitch/pkg/configuration"
	"github.com/gmcoringa/tswitch/pkg/lib"
	log "github.com/sirupsen/logrus"
)

func Init(config *configuration.Config) lib.Resolver {
	switch config.TerraformImpl {
	case "terraform":
		log.Debug("Using terraform as the terraform implementation")
		return InitTerraform()
	case "tofu":
		log.Debug("Using tofu as the terraform implementation")
		return InitTofu()
	default:
		log.Debug("Using terraform as the terraform implementation")
		return InitTerraform()
	}
}

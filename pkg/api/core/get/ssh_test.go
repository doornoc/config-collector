package get

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"testing"
)

func TestCollectConfig(t *testing.T) {
	config.GetTemplate("../../../../cmd/backend/template.json")

	s := sshStruct{Device: config.Device{
		Name:     "",
		Hostname: "",
		Port:     22,
		User:     "",
		Password: "",
		OSType:   "",
	}}
	s.accessSSHShell()
}

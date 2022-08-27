package get

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"regexp"
	"strings"
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

func TestReg(t *testing.T) {
	command := "test1 {{ command_test1 }} test2 {{ command_test2}}"
	//command := "test1 {{ "
	r := regexp.MustCompile(`\{\{.+?\}\}`)

	result := r.FindAllStringSubmatch(command, -1)
	for _, tmpCommandArray := range result {
		tmpCommand := strings.TrimSpace(tmpCommandArray[0])
		tmpCommand = strings.Replace(tmpCommand, " ", "", -1)
		tmpCommand = tmpCommand[2 : len(tmpCommand)-2]
		t.Log(tmpCommand)
	}
}

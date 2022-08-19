package debug

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"log"
)

func Deb(debugType string, deb string) {
	if config.Conf.Controller.Debug {
		log.Println(debugType, deb)
	}
}

func Err(errorType string, error error) {
	if config.Conf.Controller.Debug {
		log.Println(errorType, error)
	}
}

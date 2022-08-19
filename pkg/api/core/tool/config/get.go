package config

import (
	"fmt"
)

func GetTemplateByOSType(osType string) (OSTemplate, error) {
	for _, template := range Tpl.Templates {
		if template.OSType == osType {
			return template, nil
		}
	}

	return OSTemplate{}, fmt.Errorf("Template is not found...")
}

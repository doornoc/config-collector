package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Template struct {
	Templates []OSTemplate `json:"templates"`
}

type OSTemplate struct {
	OSType      string   `json:"os_type"`
	Commands    []string `json:"commands"`
	ConfigStart string   `json:"config_start"`
	ConfigEnd   string   `json:"config_end"`
	IgnoreLine  []string `json:"ignore_line"`
}

var Tpl Template
var TplConfigPath string

func GetTemplate(inputConfPath string) error {
	templateConfigPath := "./template.json"
	if inputConfPath != "" {
		templateConfigPath = inputConfPath
	}
	TplConfigPath = templateConfigPath
	file, err := ioutil.ReadFile(templateConfigPath)
	if err != nil {
		return err
	}
	var data Template
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	Tpl = data
	return nil
}

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Controller Controller `json:"controller"`
	Devices    []Device   `json:"devices"`
}

type Controller struct {
	TimeZone     string `json:"timezone"`
	Github       Github `json:"github"`
	TmpPath      string `json:"tmp_path"`
	SlackWebhook string `json:"slack_webhook"`
}

type Github struct {
	Repo   string `json:"repo"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Branch string `json:"branch"`
}

type Device struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	OSType   string `json:"os_type"`
}

var Conf Config

func GetConfig(inputConfPath string) error {
	configPath := "./data.json"
	if inputConfPath != "" {
		configPath = inputConfPath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Config
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	Conf = data
	return nil
}

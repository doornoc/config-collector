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
	ExecTime     uint   `json:"exec_time"`
	Debug        bool   `json:"debug"`
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
var ConfigPath string

func GetConfig(inputConfPath string) error {
	configPath := "./config.json"
	if inputConfPath != "" {
		configPath = inputConfPath
	}
	ConfigPath = configPath
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
	// when execTime == ""
	if Conf.Controller.ExecTime == 0 {
		Conf.Controller.ExecTime = 10
	}
	return nil
}

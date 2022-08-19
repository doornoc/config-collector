package get

import (
	"fmt"
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/notify"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var getConfTimer uint = 10

func CronExec() error {
	getConfTick := time.NewTicker(time.Duration(getConfTimer) * time.Second)
	getInfoTick := time.NewTicker(time.Duration(config.Conf.Controller.ExecTime) * time.Second)

	log.Printf("start for cron\n")
	for {
		select {
		case <-getConfTick.C:
			beforeNextTimer := config.Conf.Controller.ExecTime
			err := config.GetConfig(config.ConfigPath)
			if err != nil {
				log.Println(err)
				notify.NotifyErrorToSlack(err)
			}
			err = config.GetTemplate(config.TplConfigPath)
			if err != nil {
				log.Println(err)
				notify.NotifyErrorToSlack(err)
			}

			log.Printf("config timer: %d\n", config.Conf.Controller.ExecTime)
			if config.Conf.Controller.ExecTime != beforeNextTimer {
				getInfoTick = time.NewTicker(time.Duration(config.Conf.Controller.ExecTime) * time.Second)
				log.Printf("New NextTimer: %d\n", config.Conf.Controller.ExecTime)
			}
		case <-getInfoTick.C:
			err := GettingDeviceConfig()
			if err != nil {
				log.Println(err)
				notify.NotifyErrorToSlack(err)
			}
		}
	}

	return nil
}

func GettingDeviceConfig() error {
	var pushConfigs []PushConfig
	for _, device := range config.Conf.Devices {
		s := sshStruct{Device: device}
		console, err := s.accessSSHShell()
		if err != nil {
			log.Println(err)
			//return err
		}
		pushConfigs = append(pushConfigs, PushConfig{
			Name:          device.Name,
			ConfigConsole: console,
		})
	}

	err := GitPush(pushConfigs)
	if err != nil {
		return err
	}

	return nil
}

func GitPush(configs []PushConfig) error {
	if _, err := os.Stat(config.Conf.Controller.TmpPath); os.IsNotExist(err) {
		os.Mkdir(config.Conf.Controller.TmpPath, 0777)
	}

	// password authentication
	auth := &http.BasicAuth{}
	gitOption := &git.CloneOptions{}
	if config.Conf.Controller.Github.Pass != "" {
		auth = &http.BasicAuth{
			Username: config.Conf.Controller.Github.User,
			Password: config.Conf.Controller.Github.Pass,
		}
		gitOption = &git.CloneOptions{
			URL:  config.Conf.Controller.Github.Repo,
			Auth: auth,
		}
	} else {
		return fmt.Errorf("key authentication is not supported...")
	}

	var repo *git.Repository
	var err error

	repo, plainErr := git.PlainClone(config.Conf.Controller.TmpPath, false, gitOption)

	if plainErr != nil {
		repo, err = git.PlainOpen(config.Conf.Controller.TmpPath)
		log.Println(err)
		if err != nil {
			if plainErr != nil {
				log.Println("[git clone]", plainErr)
				return plainErr
			}
			log.Println("[git pull]", err)
			return err
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Println("git worktree", err)
		return err
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(config.Conf.Controller.Github.Branch),
	})
	if err != nil {
		log.Println("git checkout", err)
		return err
	}

	for _, console := range configs {
		// Create new file
		path := config.Conf.Controller.TmpPath + "/" + console.Name
		ioutil.WriteFile(path, []byte(console.ConfigConsole), 0644)
		log.Println("git add path: " + path)
		_, err = w.Add(console.Name)
		if err != nil {
			log.Println("[git add]", err)
			return err
		}
	}

	fmt.Println(w.Status())
	status, _ := w.Status()

	if status.IsClean() {
		log.Println("No need to commit")
		return nil
	}

	t := time.Now().UTC()
	fmt.Println(t)

	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	_, err = w.Commit("Updated ("+t.In(tokyo).Format(time.RFC3339)+")", &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		//RemoteName: "origin",
		Auth: auth,
	})
	if err != nil {
		return err
	}

	return nil
}

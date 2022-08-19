package get

import (
	"fmt"
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/debug"
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
				log.Println("getting config", err)
				notify.NotifyErrorToSlack(err)
			}
			err = config.GetTemplate(config.TplConfigPath)
			if err != nil {
				log.Println("getting template config", err)
				notify.NotifyErrorToSlack(err)
			}

			if config.Conf.Controller.ExecTime != beforeNextTimer {
				getInfoTick = time.NewTicker(time.Duration(config.Conf.Controller.ExecTime) * time.Second)
				log.Printf("New NextTimer: %d\n", config.Conf.Controller.ExecTime)
			}
		case <-getInfoTick.C:
			err := GettingDeviceConfig()
			if err != nil {
				log.Println("GettingDeviceConfig", err)
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
			debug.Err("[accessSSHShell]", err)
			notify.NotifyErrorToSlack(err)
		}
		pushConfigs = append(pushConfigs, PushConfig{
			Name:          device.Name,
			ConfigConsole: console,
		})
	}

	err := GitPush(pushConfigs)
	if err != nil {
		debug.Err("[GitPush]", err)
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
		debug.Deb("[Git auth]", "key authentication is not supported...")
		return fmt.Errorf("key authentication is not supported...")
	}

	var repo *git.Repository
	var err error

	repo, plainErr := git.PlainClone(config.Conf.Controller.TmpPath, false, gitOption)

	if plainErr != nil {
		repo, err = git.PlainOpen(config.Conf.Controller.TmpPath)
		if err != nil {
			if plainErr != nil {
				debug.Err("[git clone]", plainErr)
				return plainErr
			}
			debug.Err("[git pull]", err)
			return err
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		debug.Err("[git worktree]", err)
		return err
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(config.Conf.Controller.Github.Branch),
	})
	if err != nil {
		debug.Err("[git checkout]", err)
		return err
	}

	for _, console := range configs {
		// Create new file
		path := config.Conf.Controller.TmpPath + "/" + console.Name
		ioutil.WriteFile(path, []byte(console.ConfigConsole), 0644)
		debug.Deb("git add path", path)
		debug.Err("[git add]", err)

		_, err = w.Add(console.Name)
		if err != nil {
			debug.Err("[git add]", err)
			return err
		}
	}

	status, _ := w.Status()

	if status.IsClean() {
		debug.Deb("[*normal* git status] ", "No need to commit")
		return nil
	}

	t := time.Now().UTC()

	tokyo, err := time.LoadLocation(config.Conf.Controller.TimeZone)
	if err != nil {
		debug.Err("[UTC to JST]", err)
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
		debug.Err("[git commit]", err)
		return err
	}

	return nil
}

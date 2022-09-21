package get

import (
	"fmt"
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/debug"
	"github.com/doornoc/config-collector/pkg/api/core/tool/notify"
	"regexp"
	"strings"
	"time"
)

type loopStruct struct {
	stdoutUpdateTime *time.Time
	stdoutBeforeTime *time.Time
	sshSessionFinish *bool
	inCh             *chan []byte
	Device           config.Device
}

func (s *loopStruct) CommandExecLoop(command string) error {
	tickerExecCommand := time.NewTicker(time.Second * 1)
	defer tickerExecCommand.Stop()
	for {
		select {
		case <-tickerExecCommand.C:
			// compare update_time and time.Now
			now := time.Now()
			// 10:00 < 10:05
			// if now < stdoutUpdateTime
			if now.Before(s.stdoutUpdateTime.Add(time.Second * 5)) {
				continue
			}
			// if 1 minutes is break
			if now.After(s.stdoutUpdateTime.Add(time.Minute * 1)) {
				return fmt.Errorf("Command exec timeout!! (timeout)\n")
			}
			if s.stdoutBeforeTime.Equal(*s.stdoutUpdateTime) {
				return fmt.Errorf("Command exec error!! (command exec failed)\n")
			}
			// check ssh session end
			if *s.sshSessionFinish {
				return fmt.Errorf("Command exec error!! (ssh session end)\n")
			}
			// check [{{ global_variable }}] command
			result := regexp.MustCompile(`\{\{.+?\}\}`).FindAllStringSubmatch(command, -1)
			for _, tmpCommandArray := range result {
				// parse option command key
				optionCommandKey := strings.TrimSpace(tmpCommandArray[0])
				optionCommandKey = strings.Replace(optionCommandKey, " ", "", -1)
				optionCommandKey = optionCommandKey[2 : len(optionCommandKey)-2]

				value, ok := config.Conf.Options[optionCommandKey]
				if !ok {
					err := fmt.Errorf("[%s] [not found] option command (%s)", s.Device.Hostname, command)
					debug.Err("[not found] option command", err)
					notify.NotifyErrorToSlack(err)
					continue
				}

				// replace command
				command = strings.Replace(command, tmpCommandArray[0], value, -1)
			}
			*s.inCh <- []byte(command + "\n")
			*s.stdoutBeforeTime = *s.stdoutUpdateTime
			time.Sleep(1 * time.Second)
			return nil
		}
	}
}

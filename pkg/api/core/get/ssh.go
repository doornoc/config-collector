package get

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/doornoc/config-collector/pkg/api/core/tool/debug"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type sshStruct struct {
	Device config.Device `json:"device"`
}

func (s *sshStruct) accessSSHShell() (string, error) {
	consoleLog := ""

	sshConfig := &ssh.ClientConfig{
		User: s.Device.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Device.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoSKECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
			ssh.KeyAlgoSKED25519,
			ssh.KeyAlgoRSASHA256,
			ssh.KeyAlgoRSASHA512,
		},
	}

	sshConfig.KeyExchanges = append(
		sshConfig.KeyExchanges,
		"diffie-hellman-group-exchange-sha256",
		"diffie-hellman-group-exchange-sha1",
	)

	client, err := ssh.Dial("tcp", s.Device.Hostname+":"+strconv.Itoa(int(s.Device.Port)), sshConfig)
	if err != nil {
		debug.Err("[SSH Dial]", err)
		return consoleLog, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		debug.Err("[client.NewSession]", err)
		return consoleLog, err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		debug.Err("[session.StdinPipe]", err)
		return consoleLog, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		debug.Err("[session.StdoutPipe]", err)
		return consoleLog, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, 25, 80, modes)
	if err != nil {
		debug.Err("[session.RequestPty]", err)
		return consoleLog, err
	}

	err = session.Shell()
	if err != nil {
		debug.Err("[session.Shell]", err)
		return consoleLog, err
	}

	inCh := make(chan []byte)
	cancel1 := make(chan struct{})
	cancel2 := make(chan struct{})

	osTemplate, err := config.GetTemplateByOSType(s.Device.OSType)
	if err != nil {
		debug.Err("[OS Type]", err)
	}

	//stdin
	go func() {
		for {
			select {
			case b := <-inCh:
				stdin.Write(b)
				if osTemplate.InputConsole {
					consoleLog += string(b)
				}
			case <-cancel1:
				session.Close()
				return
			}
		}
	}()

	// stdout
	go func() {
		buf := make([]byte, 1000)

		for {
			var err error = nil
			for err == nil {
				select {
				case <-cancel2:
					session.Close()
					return
				default:
					n, err := stdout.Read(buf)
					consoleLog += string(buf[:n])
					if err != nil {
						debug.Err("[*normal* stdout finish]", err)
						return
					}
				}
			}
		}
	}()

	for _, command := range osTemplate.Commands {
		time.Sleep(3 * time.Second)
		inCh <- []byte(command + "\n")
	}

	time.Sleep(3 * time.Second)
	session.Close()
	// end
	close(cancel1)
	close(cancel2)

	if config.Conf.Controller.Debug {
		log.Println("==========Console==========-")
		log.Println(consoleLog)
		log.Println("====================-")
	}

	configConsole := ""
	isConfig := false
	for _, configConsoleLine := range strings.Split(consoleLog, "\n") {
		if strings.Contains(configConsoleLine, osTemplate.ConfigStart) {
			isConfig = true
		}
		if strings.Contains(configConsoleLine, osTemplate.ConfigEnd) {
			isConfig = false
		}

		if isConfig {
			isIgnoreLine := false
			for _, ignoreLine := range osTemplate.IgnoreLine {
				if strings.Contains(configConsoleLine, ignoreLine) {
					isIgnoreLine = true
					break
				}
			}
			if !isIgnoreLine {
				configConsole += configConsoleLine + "\n"
			}
		}
	}

	if config.Conf.Controller.Debug {
		log.Println("==========Config==========-")
		log.Println(configConsole)
		log.Println("====================-")
	}

	return configConsole, nil
}

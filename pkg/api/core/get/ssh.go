package get

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
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
		log.Fatal("Failed to dial: ", err)
		return consoleLog, nil
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return consoleLog, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
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
		return consoleLog, err
	}

	err = session.Shell()
	if err != nil {
		return consoleLog, err
	}

	inCh := make(chan []byte)
	cancel1 := make(chan struct{})
	cancel2 := make(chan struct{})

	//stdin
	go func() {
		for {
			select {
			case b := <-inCh:
				stdin.Write(b)
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
						log.Println("ERR")
						return
					}
				}
			}
		}
	}()

	osTemplate, err := config.GetTemplateByOSType(s.Device.OSType)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, command := range osTemplate.Commands {
		time.Sleep(3 * time.Second)
		inCh <- []byte(command + "\n")
		log.Println(command)
	}

	time.Sleep(3 * time.Second)
	err = session.Close()
	if err != nil {
		return consoleLog, err
	}
	// end
	close(cancel1)
	close(cancel2)

	log.Println("====================-")
	log.Println(consoleLog)
	log.Println("====================-")

	configConsole := ""
	isConfig := false
	for _, configConsoleLine := range strings.Split(consoleLog, "\n") {
		if strings.Contains(configConsoleLine, osTemplate.ConfigStart) {
			log.Println("test1")
			isConfig = true
		}
		if strings.Contains(configConsoleLine, osTemplate.ConfigEnd) {
			log.Println("test2")
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
	log.Println("====================-")
	log.Println(configConsole)
	log.Println("====================-")

	return configConsole, nil
}

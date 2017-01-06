package main

import (
	"os"

	"github.com/0xAX/mysql-cli/terminfo"
	"github.com/0xAX/mysql-cli/termios"
)

// Terminal structure describes current terminal where
// mysql-cli is runned.
type Terminal struct {
	inputFd  *os.File
	outputFd *os.File
	// TermCtrl provides termios capabilities (see termios(3))
	TermCtrl *termios.Termios
	// TermInfo provides TermInfo capabilities (see terminfo(5))
	TermInfo *terminfo.Terminfo
}

// InitTerm collects information about the terminal session where
// mysql-cli is runned for now and returns pointer to instance of
// the Terminal.
func InitTerm() (*Terminal, error) {
	inputFd := os.Stdin
	outputFd := os.Stdout

	termInfo, err := terminfo.ParseTerminfoFromEnv()
	if err != nil {
		return nil, err
	}

	termCtrl, err := termios.NewTermios(inputFd)
	if err != nil {
		return nil, err
	}

	termios.TcGetAttr(inputFd, termCtrl)
	terminal := &Terminal{}
	terminal.inputFd = inputFd
	terminal.outputFd = outputFd
	terminal.TermCtrl = termCtrl
	terminal.TermInfo = termInfo

	return terminal, nil
}

// IoLoop is main loop of mysql-cli process. It handles all
// input/output stuff.
func (t *Terminal) IoLoop() error {
	err := termios.NoEcho(t.inputFd, true, t.TermCtrl)
	if err != nil {
		return err
	}

	err = termios.Cbreak(t.inputFd, true, t.TermCtrl)
	if err != nil {
		return err
	}

	termios.CfMakeRaw(t.outputFd, t.TermCtrl)
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)

		if b[0] == 3 {
			break
		}

		os.Stdout.Write(b)

	}

	termios.Reset(t.outputFd, t.TermCtrl)

	return nil
}

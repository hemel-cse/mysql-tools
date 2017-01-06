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

// 	termios.CfMakeRaw(terminal.outputFd, termCtrl)
//	termios.Reset(terminal.outputFd, terminal.termCtrl)

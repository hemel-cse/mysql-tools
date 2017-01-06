package main

import (
	"os"

	"github.com/0xAX/mysql-cli/terminfo"
	"github.com/0xAX/mysql-cli/termios"
)

type Terminal struct {
	inputFd  *os.File
	outputFd *os.File
	termCtrl *termios.Termios
	termInfo *terminfo.Terminfo
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
	terminal.termCtrl = termCtrl
	terminal.termInfo = termInfo

	return terminal, nil
}

// 	termios.CfMakeRaw(terminal.outputFd, termCtrl)
//	termios.Reset(terminal.outputFd, terminal.termCtrl)

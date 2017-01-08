package main

import (
	"fmt"
	"os"

	"github.com/0xAX/mysql-tools/terminfo"
	"github.com/0xAX/mysql-tools/termios"
)

const (
	CTRL_A    = 1   // Ctrl+a
	CTRL_B    = 2   // Ctrl-b
	CTRL_C    = 3   // Ctrl-c
	CTRL_D    = 4   // Ctrl-d
	CTRL_E    = 5   // Ctrl-e
	CTRL_F    = 6   // Ctrl-f
	CTRL_H    = 8   // Ctrl-h
	TAB       = 9   // Tab
	CTRL_K    = 11  // Ctrl+k
	CTRL_L    = 12  // Ctrl+l
	ENTER     = 13  // Enter
	CTRL_N    = 14  // Ctrl-n
	CTRL_P    = 16  // Ctrl-p
	CTRL_T    = 20  // Ctrl-t
	CTRL_U    = 21  // Ctrl+u
	CTRL_W    = 23  // Ctrl+w
	ESC       = 27  // Escape
	BACKSPACE = 127 // Backspace
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
	var c []byte = make([]byte, 1)
	cmd := &SqlCommandBuffer{}

	// enable `noecho` as we will print symbols by ourself
	err := termios.NoEcho(t.outputFd, true, t.TermCtrl)
	if err != nil {
		return err
	}

	// we need in `cbreak` here to read symbol by symbol
	err = termios.Cbreak(t.outputFd, true, t.TermCtrl)
	if err != nil {
		return err
	}

	// go to raw mode
	termios.CfMakeRaw(t.outputFd, t.TermCtrl)

	//
	// start main reading loop
	//
mainloop:
	for {
		os.Stdin.Read(c)

		switch c[0] {
		case CTRL_C:
			break mainloop
		case BACKSPACE:
			if cmd.Position == 0 || cmd.Length == 0 {
				break
			}
			os.Stdout.Write([]byte(fmt.Sprintf("\r\x1b[%dC", cmd.Position-1)))
			os.Stdout.Write([]byte(" "))
			os.Stdout.Write([]byte(fmt.Sprintf("\r\x1b[%dC", cmd.Position-1)))
			cmd.Position -= 1
			cmd.Length -= 1
			break
		default:
			os.Stdout.Write(c)
			cmd.Length += 1
			cmd.Position += 1
		}
	}

	termios.Reset(t.outputFd, t.TermCtrl)

	return nil
}

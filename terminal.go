package main

import (
	"os"

	"github.com/0xAX/mysql-tools/terminfo"
	"github.com/0xAX/mysql-tools/termios"
)

const (
	ALLOC_BUFFER_STRINGS = 10
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
	UP        = 'A' // Up
	DOWN      = 'B' // Down
	RIGHT     = 'C' // Right
	LEFT      = 'D' // Left
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
	var idx int
	cmd := &SqlCommandBuffer{DisplayBuffer: make([]string, 1)}

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
			backspace(t, cmd)
			break
		case ENTER:
			idx += 1
			cmd.DisplayBuffer = append(cmd.DisplayBuffer, "")
			os.Stdout.Write([]byte("\n\r"))
			cmd.Position += 1
			cmd.Length += 1
			break
		case ESC:
			os.Stdin.Read(c)
			if c[0] == '[' {
				os.Stdin.Read(c)
				switch c[0] {
				case UP:
					/* TODO history */
					break
				case DOWN:
					/* TODO history */
					break
				case RIGHT:
					moveRight(t, cmd)
					break
				case LEFT:
					moveLeft(t, cmd)
					break
				}
			}
			break
		case TAB:
			/* TODO autocomplete */
			break
		default:
			cmd.DisplayBuffer[idx] += string(c)
			os.Stdout.Write(c)
			cmd.Length += 1
			cmd.Position += 1
			break
		}
	}

	termios.Reset(t.outputFd, t.TermCtrl)

	return nil
}

func backspace(t *Terminal, cmd *SqlCommandBuffer) {
	capability, _ := t.TermInfo.ApplyCapability("cub1")
	os.Stdout.Write([]byte(capability))
	os.Stdout.Write([]byte(" "))
	capability, _ = t.TermInfo.ApplyCapability("cub1")
	os.Stdout.Write([]byte(capability))
	cmd.Position -= 1
	cmd.Length -= 1
}

func moveRight(t *Terminal, cmd *SqlCommandBuffer) {
	if cmd.Position == cmd.Length {
		return
	}
	cmd.Position += 1
	capability, _ := t.TermInfo.ApplyCapability("cuf1")
	os.Stdout.Write([]byte(capability))
}

func moveLeft(t *Terminal, cmd *SqlCommandBuffer) {
	if cmd.Position == 0 {
		return
	}
	cmd.Position -= 1
	capability, _ := t.TermInfo.ApplyCapability("cub1")
	os.Stdout.Write([]byte(capability))
}

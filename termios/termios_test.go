// The termios functions describe a general terminal interface that
// is provided to control asynchronous communications ports.
//
// The termios library exposes all general API for a terminal controlling:
//
//  * tcgetattr
//  * tcgetattr
//  * tcsendbreak
//  * tcdrain
//  * tcflush
//  * tcflow
//  * cfmakeraw
//  * and etc.
//
// See more info at man termios(3)
package termios

import (
	"os"
	"syscall"
	"testing"
)

func TestTcgSetAttr(t *testing.T) {
	fd, err := os.OpenFile("/dev/tty", syscall.O_WRONLY|syscall.O_RDONLY, 0)
	if err != nil {
		t.Error("OpenFile /dev/tty failed")
	}

	termios, err := NewTermios(fd)
	if err != nil {
		t.Error("NewTermios failed")
	}

	result := TcGetAttr(fd, termios)
	if result != nil {
		t.Error("TcGetAttr failed")
	}

	tmpTermios := termios
	termios.c_iflag &= 0

	result = TcSetAttr(fd, TCSETS, termios)
	if result != nil {
		t.Error("TcSetAttr failed")
	}

	result = TcGetAttr(fd, termios)
	if result != nil {
		t.Error("TcGetAttr failed")
	}

	if termios.c_iflag != 0 {
		t.Error("Wrong c_iflag")
	}

	result = TcSetAttr(fd, TCSETS, tmpTermios)
	if result != nil {
		t.Error("TcSetAttr failed")
	}

	result = TcSendBreak(fd, 0)
	if result != nil {
		t.Error("tsenbreak failed")
	}

	result = TcDrain(fd)
	if result != nil {
		t.Error("TcDrain failed")
	}

	TcFlush(fd, TCIFLUSH)

	speed, err := CfGetOutputSpeed(termios)
	if err != nil || speed <= 0 {
		t.Error("CfGetOutputSpeed failed")
	}

	err = Reset(fd, termios)
	if err != nil {
		t.Error("Reset failed")
	}
}

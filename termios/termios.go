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
	"unsafe"
)

// Number of control characters
const CONTROL_CHARACTERS_NUM = 32

type cc_t byte
type speed_t uint32
type tcflag_t uint32

// Termios struct - base structure of termios package
type Termios struct {
	c_iflag     tcflag_t                     // input mode flags
	c_oflag     tcflag_t                     // output mode flags
	c_cflag     tcflag_t                     // control mode flags
	c_lflag     tcflag_t                     // local mode flags
	c_line      cc_t                         // line discipline
	c_cc        [CONTROL_CHARACTERS_NUM]cc_t // control characters
	c_ispeed    speed_t                      // input speed
	c_ospeed    speed_t                      // output speed
	origTermios *Termios                     // original Termios
}

// NewTermios initializes new termios instance and return it.
func NewTermios(fd *os.File) (*Termios, error) {
	termios := &Termios{}
	termios.origTermios = &Termios{}
	err := TcGetAttr(fd, termios.origTermios)
	if err != nil {
		return nil, err
	}
	return termios, nil
}

// TcGetAttr gets the parameters associated with the object referred by
// the give file descriptor - fd and stores them in the Termios structure
// referenced by *termios.
//
// Returns nil if everything is ok.
func TcGetAttr(fd *os.File, termios *Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}
	return nil
}

// TcSeAttr sets the parameters associated with the terminal from the
// Termios structure referred to by *termios.
//
// Returns nil if everything is ok.
func TcSetAttr(fd *os.File, optionalActions int, termios *Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(optionalActions), uintptr(unsafe.Pointer(termios)))
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}

	return nil
}

// TcSendBreak transmits a continuous stream of zero-valued bits
// for a specific duration, if the terminal is using asynchronous
// serial data transmission.
//
// If the terminal is not using asynchronous serial data transmission,
// tcsendbreak() returns without taking any action.
func TcSendBreak(fd *os.File, duration int) error {
	if duration <= 0 {
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TCSBRK), 0)
		if errno != 0 {
			return &errorTermios{errno.Error()}
		}
		return nil
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TCSBRKP), uintptr((duration+99)/100))
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}
	return nil
}

// TcDrain waits until all output written to the object referred
// to by fd has been transmitted.
func TcDrain(fd *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TCSBRK), 1)
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}

	return nil
}

// TcFlush discards data written to the object referred to by fd
// but not transmitted, or data received but not read, depending
// on the value of queueSelector which can be one of:
//
//  * TCIFLUSH - flushes data received but not read.
//  * TCOFLUSH - flushes data written but not transmitted.
//  * TCIOFLUSH - flushes both data received but not read,
// and data written but not transmitted.
//
func TcFlush(fd *os.File, queueSelector int) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TCFLSH), uintptr(queueSelector))
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}
	return nil
}

// TcFlow suspends transmission or reception of data on the object
// referred to by fd, depending on the value of action:
//
//  * TCOOFF suspends output.
//  * TCOON  restarts suspended output.
//  * TCIOFF transmits a STOP character, which stops the terminal
// device from transmitting data to the system.
//  * TCION  transmits a START character, which starts the terminal
// device transmitting data to the system.
//
func TcFlow(fd *os.File, action int) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TCXONC), uintptr(action))
	if errno != 0 {
		return &errorTermios{errno.Error()}
	}
	return nil
}

// CfGetOutputSpeed returns the output baud rate stored in the Termios
// structure pointed to by *termios.
func CfGetOutputSpeed(termios *Termios) (speed_t, error) {
	if termios == nil {
		return 0, &errorTermios{"termios must be initialized"}
	}
	ospeed := termios.c_cflag & (CBAUD | CBAUDEX)
	return speed_t(ospeed), nil
}

// CfGetInputSpeed returns the input baud rate stored in the termios structure.
//
// There is no difference between input and output
// speed for Linux.
func CfGetInputSpeed(termios *Termios) (speed_t, error) {
	return CfGetOutputSpeed(termios)
}

// Reset function resets current state of Termios to the
// old value pointed by the termios.
func Reset(fd *os.File, termios *Termios) error {
	return TcSetAttr(fd, TCSETSF, termios.origTermios)
}

// CfMakeRaw sets the terminal to something like the "raw" mode of the old
// Version 7 terminal driver: input is available character by character,
// echoing is disabled, and  all  special processing of terminal input and
// output characters is disabled.
func CfMakeRaw(fd *os.File, termios *Termios) error {
	if termios == nil {
		return &errorTermios{"termios must be initialized"}
	}
	termios.c_iflag &^= (IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | ICRNL | IXON)
	termios.c_oflag &^= OPOST
	termios.c_lflag &^= (ECHO | ECHONL | ICANON | ISIG | IEXTEN)
	termios.c_cflag &^= (CSIZE | PARENB)
	termios.c_cflag |= CS8
	termios.c_cc[VMIN] = 1
	termios.c_cc[VTIME] = 0

	return TcSetAttr(fd, TCSETSF, termios)
}

// NoEcho controls whether input is immediately re-echoed as output.
// If enable is true, then echo will be disabled.
func NoEcho(fd *os.File, enable bool, termios *Termios) error {
	if termios == nil {
		return &errorTermios{"termios must be initialized"}
	}

	if enable == true {
		termios.c_lflag ^= ECHO
	} else {
		termios.c_lflag |= ECHO
	}

	return TcSetAttr(fd, TCSETSF, termios)
}

// Cbreak enables/disables `canonical` mode. If `canonical` mode is disabled,
// the input will be returned to program immidiately.
func Cbreak(fd *os.File, enable bool, termios *Termios) error {
	if termios == nil {
		return &errorTermios{"termios must be initialized"}
	}

	if enable == true {
		termios.c_lflag ^= ICANON
	} else {
		termios.c_lflag |= ICANON
	}

	return TcSetAttr(fd, TCSETSF, termios)
}

// CfSetOutputSpeed sets the output baud rate stored in the Termios structure
// pointed to by Termios to speed, which must be one of the BAUD*
// constants (see above)
func CfSetOutputSpeed(termios *Termios, speed speed_t) error {
	if termios == nil {
		return &errorTermios{"termios must be initialized"}
	}
	if termios.c_ispeed&^CBAUD != 0 && termios.c_ispeed < B57600 || termios.c_ispeed > MAX_BAUD {
		return &errorTermios{"BAUD EINVAL"}
	}
	termios.c_ospeed = speed
	termios.c_cflag &^= (CBAUD | CBAUDEX)
	termios.c_cflag |= tcflag_t(speed)

	return nil
}

// CfSetInputSpeed sets the input baud rate stored in the Termios structure
// to speed, which must be specified as one of the Bnnn constants
// listed above for cfsetospeed().
func CfSetInputSpeed(termios *Termios, speed speed_t) error {
	if termios == nil {
		return &errorTermios{"termios must be initialized"}
	}
	if termios.c_ispeed&^CBAUD != 0 && termios.c_ispeed < B57600 || termios.c_ispeed > MAX_BAUD {
		return &errorTermios{"BAUD EINVAL"}
	}
	termios.c_ispeed = speed

	// if the given speed is 0, input speed will be
	// the same that output speed
	if speed == 0 {
		termios.c_ispeed = termios.c_ospeed
	} else {
		termios.c_cflag &^= (CBAUD | CBAUDEX)
		termios.c_cflag |= tcflag_t(speed)
	}
	return nil
}

// GetInputFlags returns input flags from Termios
func (t *Termios) GetInputFlags() (tcflag_t, error) {
	return t.c_iflag, nil
}

// GetOutputFlags returns output flags from Termios
func (t *Termios) GetOutputFlags() (tcflag_t, error) {
	return t.c_oflag, nil
}

// GetControlFlags returns control flags from Termios
func (t *Termios) GetControlFlags() (tcflag_t, error) {
	return t.c_cflag, nil
}

// GetLocalModesFlags return local modes flags from Termios
func (t *Termios) GetLocalModesFlags() (tcflag_t, error) {
	return t.c_lflag, nil
}

// GetLineDiscipline returns line control flags
func (t *Termios) GetLineDiscipline() (cc_t, error) {
	return t.c_line, nil
}

// GetControlCharacters returns control characters
func (t *Termios) GetControlCharacters() ([32]cc_t, error) {
	return t.c_cc, nil
}

// GetInputSpeed returns input speed
func (t *Termios) GetInputSpeed() (speed_t, error) {
	return t.c_ispeed, nil
}

// GetInputOutputSpeed returns output speed
func (t *Termios) GetInputOutputSpeed() (speed_t, error) {
	return t.c_ospeed, nil
}

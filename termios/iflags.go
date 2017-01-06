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

// termios c_iflag flag constants
const (
	IGNBRK  = 0000001 // Ignore BREAK.
	BRKINT  = 0000002 // Causes the input and output queues to be flushed.
	IGNPAR  = 0000004 // Ignore framing errors and parity errors.
	PARMRK  = 0000010 // If IGNPAR is not set, prefix a character with a parity error or framing error with \377 \0. In other way only \0.
	INPCK   = 0000020 // Enable input parity checking.
	ISTRIP  = 0000040 // Strip off eighth bit.
	INLCR   = 0000100 // Translate NL to CR on input.
	IGNCR   = 0000200 // Ignore carriage return on input.
	ICRNL   = 0000400 // Translate carriage return to newline on input (unless IGNCR is set).
	IUCLC   = 0001000 // Map uppercase characters to lowercase on input.
	IXON    = 0002000 // Enable XON/XOFF flow control on output.
	IXANY   = 0004000 // Typing any character will restart stopped output.
	IXOFF   = 0010000 // Enable XON/XOFF flow control on input.
	IMAXBEL = 0020000 // Ring bell when input queue is full.
	IUTF8   = 0040000 // Input is UTF8; this allows character-erase to be correctly performed in cooked mode.
)

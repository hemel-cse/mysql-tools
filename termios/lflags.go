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

// termios c_lflag flag constants
const (
	ISIG    = 0000001 // When any of the characters INTR, QUIT, SUSP, or DSUSP are received, generate the corresponding signal.
	ICANON  = 0000002 // Enable canonical mode.
	XCASE   = 0000004 // If ICANON is also set, terminal is uppercase only.
	ECHO    = 0000010 // Echo input characters.
	ECHOE   = 0000020 // If ICANON is also set, the ERASE character erases the preceding input character, and WERASE erases the preceding word.
	ECHOK   = 0000040 // If ICANON is also set, the KILL character erases the current line.
	ECHONL  = 0000100 // If ICANON is also set, echo the NL character even if ECHO is not set.
	NOFLSH  = 0000200 // Disable flushing the input and output queues when generating signals for the INT, QUIT, and SUSP characters.
	TOSTOP  = 0000400 // Send the SIGTTOU signal to the process group of a background process which tries to write to its controlling terminal.
	ECHOCTL = 0001000 // If ECHO is also set, terminal special characters other than TAB, NL, START, and STOP are echoed as ^X,
	// where X is the character with ASCII code 0x40 greater than the special character.
	ECHOPRT = 0002000 // If ICANON and ECHO are also set, characters are printed as they are being erased.
	ECHOKE  = 0004000 // If ICANON is also set, KILL is echoed by erasing each character on the line, as specified by ECHOE and ECHOPRT.
	FLUSHO  = 0010000 // Output is being flushed.
	PENDIN  = 0040000 // All characters in the input queue are reprinted when the next character is read.
	IEXTEN  = 0100000 // Enable implementation-defined input processing.
)

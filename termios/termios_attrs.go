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

// termios tcsetattr, tcflush and other API related options
const (
	TCSETS  = 0x5402 // the change occurs immediately.
	TCSETSW = 0x5403 // the change occurs after all output written to fd has been transmitted.
	// This option should be used when changing parameters that affect output.
	TCSETSF = 0x5404 // the change occurs after all output written to the object
	// referred by fd has been transmitted, and all input that has
	// been received but not read will be discarded before the change is made.
	TCSBRK    = 0x5409 // Send a break.
	TCSBRKP   = 0x5425 // So-called "POSIX version" of TCSBRK
	TCIFLUSH  = 0      // flushes data received but not read.
	TCOFLUSH  = 1      // flushes data written but not transmitted.
	TCIOFLUSH = 2      // flushes both data received but not read, and data written but not transmitted.
	TCFLSH    = 0x540B // Equivalent to tcflush(fd, arg).
	TCXONC    = 0x540A // Equivalent to tcflow(fd, arg).
	TIOCEXCL  = 0x540C // Put the terminal into exclusive mode. No further open(2) operations on the terminal are permitted.
	TCOOFF    = 0      // suspends output.
	TCOON     = 1      // restarts suspended output.
	TCIOFF    = 2      // transmits a STOP character, which stops the terminal device from transmitting data to the system.
	TCION     = 3      // transmits a START character, which starts the terminal device transmitting data to the system
)

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

// c_cflag flag constants
const (
	B0       = 0000000 // From B0 to B4000000 line speed
	B50      = 0000001
	B75      = 0000002
	B110     = 0000003
	B134     = 0000004
	B150     = 0000005
	B200     = 0000006
	B300     = 0000007
	B600     = 0000010
	B1200    = 0000011
	B1800    = 0000012
	B2400    = 0000013
	B4800    = 0000014
	B9600    = 0000015
	B19200   = 0000016
	B38400   = 0000017
	B57600   = 0010001
	B115200  = 0010002
	B230400  = 0010003
	B460800  = 0010004
	B500000  = 0010005
	B921600  = 0010007
	B1000000 = 0010010
	B1152000 = 0010011
	B1500000 = 0010012
	B2000000 = 0010013
	B2500000 = 0010014
	B3000000 = 0010015
	B3500000 = 0010016
	B4000000 = 0010017
	MAX_BAUD = B4000000
	CSIZE    = 0000060 // Character size mask.  Values are CS5, CS6, CS7, or CS8.
	CS5      = 0000000
	CS6      = 0000020
	CS7      = 0000040
	CS8      = 0000060
	CSTOPB   = 0000100 // Set two stop bits, rather than one.
	CREAD    = 0000200 // Enable receiver.
	PARENB   = 0000400 // Enable parity generation on output and parity checking for input.
	PARODD   = 0001000 // If set, then parity for input and output is odd; otherwise even parity is used.
	HUPCL    = 0002000 // Lower modem control lines after last process closes the device (hang up).
	CLOCAL   = 0004000 // Ignore modem control lines.
	CBAUDEX  = 0010000 // Extra baud speed mask (1 bit), included in CBAUD.
	CBAUD    = 0010017 // Baud speed mask (4+1 bits)
)

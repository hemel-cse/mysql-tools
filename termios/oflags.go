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

// termios c_oflag constants
const (
	OPOST  = 0000001 // Enable implementation-defined output processing.
	OLCUC  = 0000002 // Map lowercase characters to uppercase on output.
	ONLCR  = 0000004 // Map NL to CR-NL on output.
	OCRNL  = 0000010 // Map CR to NL on output.
	ONOCR  = 0000020 // Don't output CR at column 0.
	ONLRET = 0000040 // Don't output CR.
	OFILL  = 0000100 // Send fill characters for a delay, rather than using a timed delay.
	OFDEL  = 0000200 // Fill character is ASCII DEL (0177).
	NLDLY  = 0000400 // Newline delay mask.
	NL0    = 0000000 //
	NL1    = 0000400 //
	CRDLY  = 0003000 // Carriage return delay mask.
	CR0    = 0000000
	CR1    = 0001000
	CR2    = 0002000
	CR3    = 0003000
	TABDLY = 0014000 // Horizontal  tab  delay  mask.
	TAB0   = 0000000
	TAB1   = 0004000
	TAB2   = 0010000
	TAB3   = 0014000
	BSDLY  = 0020000 // Backspace delay mask.
	BS0    = 0000000
	BS1    = 0020000
	FFDLY  = 0100000 // Form feed delay mask.
	FF0    = 0000000
	FF1    = 0100000
)

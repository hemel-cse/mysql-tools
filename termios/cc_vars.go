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

const (
	VINTR    = 0  // Interrupt character [ISIG].
	VQUIT    = 1  // Quit character [ISIG].
	VERASE   = 2  // Erase character [ICANON].
	VKILL    = 3  // Kill-line character [ICANON].
	VMIN     = 4  // Minimum number of bytes read at once [!ICANON].
	VTIME    = 5  // Time-out value (tenths of a second) [!ICANON].
	VEOL2    = 6  // Second EOL character [ICANON].
	VSWTCH   = 7  // Used in System V to switch shells in shell layers.
	VSTART   = 8  // Start (X-ON) character [IXON, IXOFF].
	VSTOP    = 9  // Stop (X-OFF) character [IXON, IXOFF].
	VSUSP    = 10 // Suspend character [ISIG].
	VDSUSP   = 11 // Delayed suspend character, send SIGTSTP signal when the character is read by the user program .
	VREPRINT = 12 // Reprint-line character [ICANON].
	VDISCARD = 13 // Start/stop discarding pending output.
	VWERASE  = 14 // Word-erase character [ICANON].
	VLNEXT   = 15 // Literal-next character [IEXTEN].
	VEOF     = 16 // End-of-file character [ICANON].
	VEOL     = 17 // End-of-line character [ICANON].
)

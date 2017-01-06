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

// errorTermios is an implementation of error interface which
// provides error message for the API of termios package. Mostly
// used as wrapper for Errno.
type errorTermios struct {
	errMsg string
}

func (err errorTermios) Error() string {
	return err.errMsg
}

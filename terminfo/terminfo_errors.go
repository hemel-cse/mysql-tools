// terminfo package provides simple parser for compiled terminfo fields
// and additionally it provides trivial API for applying of capability
// of a terminal.
package terminfo

// errorTerminfo is an implementation of error interface which
// provides error message for the ParseTerminfo().
type errorTerminfo struct {
	errMsg string
}

func (err errorTerminfo) Error() string {
	return err.errMsg
}

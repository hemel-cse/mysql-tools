package main

// SqlCommnadBuffer provides interface to current sql command typed
// by an user.
type SqlCommandBuffer struct {
	// Length of current sql command
	Length uint64
	// Current position in sql command
	Position uint64
	// Current buffer with colors and etc. This one will not be
	// sent to MySql server. It is only for displaying.
	DisplayBuffer []string
	// Current line index in DisplayBuffer
	DisplayBufferIndex uint64
}

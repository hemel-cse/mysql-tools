package main

// SqlCommnadBuffer provides interface to current sql command typed
// by an user.
type SqlCommandBuffer struct {
	// Length of current sql command
	Length uint64
	// Current position in sql command
	Position uint64
}

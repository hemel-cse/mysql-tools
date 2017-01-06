package main

var (
	// Databases is list of databases avaliable for current user
	Databases []string = []string{"mysql"}
)

// DatabaseList provides list of databases which are avaliable for
// current user
func DatabaseList() func(line string) []string {
	return func(line string) []string {
		return Databases
	}
}

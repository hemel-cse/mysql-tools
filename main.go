package main

import "flag"

func main() {
	/* parse mysql-cli flags (defined in flags.go) */
	flag.Parse()

	/* TODO connect to database */
	/* TODO parse configuration */

	/* initialize terminal and collect information about current terminal session  */
	terminal, err := InitTerm()
	if err != nil {
		panic(err)
	}

	/* start main loop */
	terminal.IoLoop()
}

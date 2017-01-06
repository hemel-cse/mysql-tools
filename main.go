package main

import "flag"

func main() {
	/* parse mysql-cli flags (defined in flags.go) */
	flag.Parse()

	/*  */
	InitTerm()

	/* TODO connect to database */
	/* TODO parse configuration */
}

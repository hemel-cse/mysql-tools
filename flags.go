package main

import "flag"

var (
	bind_address = flag.String("bind-address", "", `On a computer having multiple network interfaces, 
        use this option to select which interface to use 
        for connecting to the MySQL server.`)

	db = flag.String("database,-D", "mysql", `The database to use. This is useful primarily in 
        an option file.`)

	user = flag.String("user,-u", "root", `The MySQL user name to use when connecting to the 
        server.`)

	password = flag.String("password,-p", "", `The password to use when connecting to the server. 
        If you use the short option form (-p), you cannot 
        have a space between the option and the password. 
        If you omit the password value following the 
        --password or -p option on the command line, mysql 
        prompts for one.`)

	history_file = flag.String("history", "", `The fill will be used as history for users
	commands.`)
)

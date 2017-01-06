package main

import "github.com/chzyer/readline"

var completer = readline.NewPrefixCompleter(
	readline.PcItem("USE",
		readline.PcItemDynamic(DatabaseList()),
	),
)

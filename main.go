// These examples demonstrate more intricate uses of the flag package.
package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	var cmds []string

	flag.Parse()

	/* TODO connect to database */
	/* TODO parse configuration */

	/* configure initial readline */
	l, err := readline.NewEx(&readline.Config{
		Prompt:                 "\033[31m»\033[0m ",
		HistoryFile:            "/tmp/readline.tmp",
		AutoComplete:           completer,
		InterruptPrompt:        "^C",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	/* main loop */
	for {
		line, err := l.Readline()
		if err != nil {
			break
		}

		/* clean our command line from spaces */
		line = strings.TrimSpace(line)

		/* skip empty lines */
		if len(line) == 0 {
			continue
		}

		cmds = append(cmds, line)

		/* we didn't get ';' yet, so let's continue to read */
		if !strings.HasSuffix(line, ";") {
			continue
		}

		/* join into one line mutli-line commands */
		cmd := strings.Join(cmds, " ")

		/* update history */
		l.SaveHistory(cmd)

		/* start to execute db query */
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
}

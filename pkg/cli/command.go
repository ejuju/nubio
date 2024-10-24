package cli

import "slices"

type Command struct {
	Keyword     string
	Usage       string
	Description string
	Aliases     []string
	Do          func(args ...string) (exitcode int)
}

func Index(keyword string, commands []*Command) *Command {
	for _, cmd := range commands {
		if keyword == cmd.Keyword || slices.Contains(cmd.Aliases, keyword) {
			return cmd
		}
	}
	return nil
}

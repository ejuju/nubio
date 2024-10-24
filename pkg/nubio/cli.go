package nubio

import (
	"fmt"
	"log"
	"os"

	"github.com/ejuju/nubio/pkg/cli"
)

func Run(args ...string) (exitcode int) {
	if len(args) == 0 {
		return commandRunServer.Do()
	}

	cmd := cli.Index(args[0], commands)
	if cmd == nil {
		fmt.Printf("Unknown command: %q.\n", args[0])
		fmt.Printf("Use %q to list available commands.\n", commandHelp.Keyword)
		return 1
	}
	return cmd.Do(args[1:]...)
}

var commands = []*cli.Command{
	commandRunServer,
	commandRunSSG,
	commandGeneratePDF,
	commandCheckProfile,
	commandCheckServer,
}

// Prepend help command.
func init() { commands = append([]*cli.Command{commandHelp}, commands...) }

var commandHelp = &cli.Command{
	Keyword:     "help",
	Aliases:     []string{"--help", "?", "menu"},
	Description: "Print available commands.",
	Do: func(args ...string) (exitcode int) {
		fmt.Printf("Available commands:\n")
		for _, cmd := range commands {
			fmt.Printf("- \x1b[33m%-14s\x1b[0m %s ", cmd.Keyword, cmd.Description)
			if cmd.Usage != "" {
				fmt.Printf("\x1b[30mExample: %q.\x1b[0m", cmd.Usage)
			}
			fmt.Print("\n")
		}
		return 0
	},
}

var commandRunServer = &cli.Command{
	Keyword:     "run",
	Usage:       "run $PATH_TO_SERVER_CONF",
	Description: "Run as HTTP(S) server.",
	Do:          RunServer,
}

var commandRunSSG = &cli.Command{
	Keyword:     "ssg",
	Usage:       "ssg $PATH_TO_PROFILE_CONF $PATH_TO_OUTPUT_DIR",
	Description: "Generate static website files.",
	Do:          RunSSG,
}

var commandGeneratePDF = &cli.Command{
	Keyword:     "pdf",
	Usage:       "pdf $PATH_TO_PROFILE_CONF $PATH_TO_OUTPUT_DIR",
	Description: "Generate PDF export.",
	Do: func(args ...string) (exitcode int) {
		if len(args) < 2 {
			log.Println("missing argument(s): /path/to/profile.json and /path/to/output.pdf")
		}
		in := args[0]
		out := args[1]

		// Load and check profile.json.
		p := &Profile{}
		err := loadJSONFile(in, p)
		if err != nil {
			log.Printf("load profile: %s", err)
			return 1
		}
		errs := p.Check()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("check profile: %s", err)
			}
			return 1
		}

		// Encode and write PDF.
		f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("open output file: %s", err)
			return 1
		}
		defer f.Close()
		err = ExportPDF(f, p)
		if err != nil {
			log.Printf("encode and write PDF: %s", err)
			return 1
		}

		log.Printf("PDF written to %s", out)
		return 0
	},
}

var commandCheckProfile = &cli.Command{
	Keyword:     "check-profile",
	Usage:       "check-profile $PATH_TO_PROFILE_CONF",
	Description: "Check a \"profile.json\" file.",
	Do: func(args ...string) (exitcode int) {
		// Load profile.json.
		path := "profile.json"
		if len(args) > 0 {
			path = args[0]
		}
		log.Printf("Checking file: %s", path)
		p := &Profile{}
		err := loadJSONFile(path, p)
		if err != nil {
			return 1
		}
		errs := p.Check()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error: %s", err)
			}
			return 1
		}

		log.Printf("All good!")
		return 0
	},
}

var commandCheckServer = &cli.Command{
	Keyword:     "check-server",
	Usage:       "check-server $PATH_TO_SERVER_CONF",
	Description: "Check a \"server.json\" file.",
	Do: func(args ...string) (exitcode int) {
		// Load server.json.
		path := "server.json"
		if len(args) > 0 {
			path = args[0]
		}
		log.Printf("Checking file: %s", path)
		p := &Config{}
		err := loadJSONFile(path, p)
		if err != nil {
			return 1
		}
		errs := p.Check()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error: %s", err)
			}
			return 1
		}

		log.Printf("All good!")
		return 0
	},
}

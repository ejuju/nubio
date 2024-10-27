package nubio

import (
	"fmt"
	"log"
	"os"

	"github.com/ejuju/nubio/pkg/cli"
)

const version = "v0"

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
	commandVersion,
	commandRunServer,
	commandRunSSG,
	commandExport,
	commandCheckConfig,
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

var commandVersion = &cli.Command{
	Keyword:     "version",
	Aliases:     []string{"v", "-v", "--v", "-version", "--version"},
	Description: "Print the version of this executable.",
	Do: func(args ...string) (exitcode int) {
		fmt.Printf("%s\n", version)
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
	Usage:       "ssg $PATH_TO_CONFIG $PATH_TO_OUTPUT_DIR",
	Description: "Generate static website files.",
	Do:          RunSSG,
}

var commandExport = &cli.Command{
	Keyword:     "export",
	Usage:       "export $FORMAT $CONFIG_PATH $OUTPUT_PATH",
	Description: "Export to file.",
	Do: func(args ...string) (exitcode int) {
		if len(args) < 3 {
			log.Println("missing argument(s): format, config_path, output_path")
			return 1
		}
		format := ExportFormat(args[0])
		in := args[1]
		out := args[2]

		// Load and check config.json.
		conf, err := LoadConfig(in)
		if err != nil {
			log.Printf("load config: %s", err)
			return 1
		}
		errs := conf.Check()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("check config: %s", err)
			}
			return 1
		}

		// Encode and write.
		var exporter ExportFunc
		switch format {
		default:
			log.Printf("unknown export format: %q", format)
			return 1
		case ExportTypeHTML:
			exporter = ExportHTML
		case ExportTypePDF:
			exporter = ExportPDF
		case ExportTypeJSON:
			exporter = ExportJSON
		case ExportTypeTXT:
			exporter = ExportText
		case ExportTypeMD:
			exporter = ExportMarkdown
		}
		f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("open output file: %s", err)
			return 1
		}
		defer f.Close()
		err = exporter(f, conf)
		if err != nil {
			log.Printf("encode and write: %s", err)
			return 1
		}

		log.Printf("wrote export to %s", out)
		return 0
	},
}

var commandCheckConfig = &cli.Command{
	Keyword:     "check",
	Aliases:     []string{"check-config"},
	Usage:       "check $PATH_TO_CONFIG",
	Description: "Check a \"config.json\" file.",
	Do: func(args ...string) (exitcode int) {
		// Load config.json.
		path := "config.json"
		if len(args) > 0 {
			path = args[0]
		}
		log.Printf("Checking file: %s", path)
		conf, err := LoadConfig(path)
		if err != nil {
			return 1
		}
		errs := conf.Check()
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

package nubio

import (
	"log"
	"os"
)

func Run(args ...string) (exitcode int) {
	if len(args) == 0 {
		return RunServer()
	}

	var do func(args ...string) (exitcode int)
	switch cmd := args[0]; cmd {
	default:
		log.Fatalf("unknown command %q", cmd)
		return 1
	case "run":
		do = RunServer
	case "ssg":
		do = RunSSG
	case "pdf":
		do = RunPDF
	case "check-profile":
		do = RunCheckProfile
	case "check-server-config":
		do = RunCheckServer
	}

	return do(args[1:]...)
}

func RunPDF(args ...string) (exitcode int) {
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
}

func RunCheckProfile(args ...string) (exitcode int) {
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
}

func RunCheckServer(args ...string) (exitcode int) {
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
}

package nubio

import "log"

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
	case "check-profile-config":
		do = RunCheckProfile
	case "check-server-config":
		do = RunCheckServer
	}

	return do(args[1:]...)
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

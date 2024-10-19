package nubio

import "log"

func Run(args ...string) (exitcode int) {
	if len(args) == 0 {
		return RunServer()
	}

	switch cmd := args[0]; cmd {
	default:
		log.Fatalf("unknown command %q", cmd)
		return 1
	case "run":
		return RunServer(args[1:]...)
	case "ssg":
		return RunSSG(args[1:]...)
	}
}

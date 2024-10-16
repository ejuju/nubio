package nuage

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func Run(args ...string) (exitcode int) {
	slogh := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(slogh)
	logger.Debug("logger ready")

	go http.ListenAndServe(":8080", http.NotFoundHandler())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)
	<-interrupt

	logger.Debug("exiting gracefully")
	return 0
}

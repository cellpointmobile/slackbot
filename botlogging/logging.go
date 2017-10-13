package botlogging

import (
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("")

var format = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

func SetupLogging(level logging.Level) {
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")

	logging.SetBackend(backendLeveled)

	log.Info("Initialized logging with level: " + level.String() )
}

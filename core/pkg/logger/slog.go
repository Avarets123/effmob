package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/MatusOllah/slogcolor"
)

var logger *slog.Logger

func init() {
	logger = newSlog()
}

func GetLogger() *slog.Logger {
	return logger
}

func newSlog() *slog.Logger {

	return slog.New(slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
		SrcFileMode: slogcolor.Nop,
		Level:       slog.LevelInfo,
		NoColor:     false,
		TimeFormat:  time.DateTime,
	}))
}

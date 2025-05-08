package lg

import (
	"log/slog"
	"os"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case Local:
		//выводить все логи в текстовом виде в консоль (os.Stdout)
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Dev:
		//выводить все логи в Json виде в консоль (os.Stdout) в деве, полезно для
		//различных сервисов, прометеус, кибана и тд
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case Prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("unknown env")
	}

	return log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

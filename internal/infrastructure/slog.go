package infrastructure

import (
	"log/slog"
	"os"
)

func InitSlog() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, value slog.Attr) slog.Attr {
			return value
		},
	}).WithAttrs([]slog.Attr{
		slog.String("service", "sfa-service"),
		slog.String("with-release", "v1.0.0"),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

package utils

import (
	"log/slog"
	"os"
)

var Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

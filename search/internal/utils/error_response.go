package utils

import "log/slog"

func ErrorResponse(err error, message string) map[string]string {
	slog.Error(message, "error", err)
	return map[string]string{"error": message}
}

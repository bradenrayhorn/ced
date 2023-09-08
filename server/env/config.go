package env

import (
	"os"
	"strings"

	"github.com/bradenrayhorn/ced/server/ced"
)

func LoadConfig() ced.Config {
	return ced.Config{
		PrettyLog:     isTrue(os.Getenv("PRETTY_LOG")),
		DbPath:        GetDbPath(),
		HttpPort:      os.Getenv("HTTP_PORT"),
		AllowedOrigin: os.Getenv("ORIGIN"),
	}
}

func isTrue(string string) bool {
	string = strings.ToLower(string)
	return string == "yes" || string == "true"
}

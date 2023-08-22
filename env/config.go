package env

import (
	"os"
	"strings"

	"github.com/bradenrayhorn/ced/ced"
)

func LoadConfig() ced.Config {
	return ced.Config{
		PrettyLog: isTrue(os.Getenv("PRETTY_LOG")),
		DbPath:    os.Getenv("DB_PATH"),
		HttpPort:  os.Getenv("HTTP_PORT"),
	}
}

func isTrue(string string) bool {
	string = strings.ToLower(string)
	return string == "yes" || string == "true"
}

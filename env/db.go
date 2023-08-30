package env

import "os"

func GetDbPath() string {
	return os.Getenv("DB_PATH")
}

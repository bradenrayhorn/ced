package env

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestIsTrue(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected bool
	}{
		{"yes", "yes", true},
		{"true", "true", true},
		{"empty", "", false},
		{"ignores case", "tRUe", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			res := isTrue(test.input)

			is.Equal(test.expected, res)
		})
	}
}

func TestReadsFromEnv(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("DB_PATH", "path.db"))
	is.NoErr(os.Setenv("PRETTY_LOG", "true"))

	config := LoadConfig()

	is.Equal(config.DbPath, "path.db")
	is.Equal(config.PrettyLog, true)
}

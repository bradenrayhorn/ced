package ced

import "errors"

type Config struct {
	PrettyLog bool
	DbPath    string
	HttpPort  string
	// CORS allowed origins
	AllowedOrigin string

	// If set, the client IP will be read from this header.
	TrustedClientIPHeader string
}

func (c Config) Validate() error {
	err := ValidateFields(
		Field("DbPath", Required(ValidatableString(c.DbPath))),
		Field("HttpPort", Required(ValidatableString(c.HttpPort))),
	)

	var cedError Error
	if errors.As(err, &cedError) {
		_, msg := cedError.CedError()
		return errors.New(msg)
	}

	return err
}

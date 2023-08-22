package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bradenrayhorn/ced/ced"
)

func decodeRequest(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		// EOF means no body was sent and can be treated as non-error (further validation should likely fail)
		if errors.Is(err, io.EOF) {
			return nil
		}

		// Some JSON syntax errors are returned as an io.ErrUnexpectedEOF rather than json.SyntaxError.
		// See https://github.com/golang/go/issues/25956 for reference.
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return ced.NewError(ced.EUNPROCESSABLE, "Invalid JSON provided.")
		}

		// JSON syntax error
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			return ced.NewError(ced.EUNPROCESSABLE, syntaxError.Error())
		}

		// JSON type error
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			return ced.NewError(
				ced.EUNPROCESSABLE,
				fmt.Sprintf(
					"Invalid data `%s` for `%s` field.",
					unmarshalTypeError.Value,
					unmarshalTypeError.Field,
				),
			)
		}

		return err
	}

	return nil
}

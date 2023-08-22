package ced

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/segmentio/ksuid"
)

func TestFromString(t *testing.T) {
	is := is.New(t)

	t.Run("parses id", func(t *testing.T) {
		id, err := IDFromString("0ujsswThIGTUYm2K8FjOOfXtY1K")
		is.NoErr(err)
		is.Equal(id.String(), "0ujsswThIGTUYm2K8FjOOfXtY1K")
	})

	t.Run("handles parse error", func(t *testing.T) {
		_, err := IDFromString("jibberish")
		is.True(err != nil)
	})

	t.Run("handles empty string", func(t *testing.T) {
		_, err := IDFromString("jibberish")
		is.True(err != nil)
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Run("empty id", func(t *testing.T) {
		is := is.New(t)

		res, err := json.Marshal(ID(ksuid.Nil))
		is.NoErr(err)
		is.Equal(string(res), "null")
	})

	t.Run("non-empty id", func(t *testing.T) {
		is := is.New(t)

		id, _ := IDFromString("2UK7TWjLNyVeTcBfkzt36InToZr")
		res, err := json.Marshal(id)
		is.NoErr(err)
		is.Equal(string(res), `"2UK7TWjLNyVeTcBfkzt36InToZr"`)
	})
}

func TestUnmarshalJSON(t *testing.T) {
	t.Run("can unmarshal", func(t *testing.T) {
		is := is.New(t)

		var id ID
		err := json.Unmarshal([]byte(`"2UK7bBkOso71dHZRLAcglrnRqqp"`), &id)
		is.NoErr(err)
		is.Equal(id.String(), "2UK7bBkOso71dHZRLAcglrnRqqp")
	})

	t.Run("handles invalid id", func(t *testing.T) {
		is := is.New(t)

		var id ID
		err := json.Unmarshal([]byte(`"&*1"`), &id)
		is.Equal(err.Error(), "invalid id: &*1")
	})

	t.Run("handles json parse error", func(t *testing.T) {
		is := is.New(t)

		var id ID
		err := json.Unmarshal([]byte(`"13`), &id)
		is.Equal(err.Error(), "unexpected end of JSON input")
	})
}

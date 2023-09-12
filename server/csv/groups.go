package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/bradenrayhorn/ced/server/ced"
)

func ParseGroupImport(reader io.Reader) ([]ced.GroupImport, error) {
	// csv has format, no header line:
	// {name}, {max_attendees}, {search_hints}

	r := csv.NewReader(reader)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = 3

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv: %w", err)
	}

	res := make([]ced.GroupImport, len(records))
	for i, record := range records {
		g, err := parseRecord(record)
		if err != nil {
			return nil, fmt.Errorf("failed to parse record on line %d: %w", i+1, err)
		}

		res[i] = g
	}

	return res, nil
}

func parseRecord(record []string) (ced.GroupImport, error) {
	maxAttendees, err := strconv.ParseUint(record[1], 10, 8)
	if err != nil {
		return ced.GroupImport{}, fmt.Errorf("invalid max_attendees: %s", record[1])
	}

	return ced.GroupImport{
		Name:         ced.Name(record[0]),
		MaxAttendees: uint8(maxAttendees),
		SearchHints:  record[2],
	}, nil
}

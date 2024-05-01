package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/bradenrayhorn/ced/server/ced"
)

// import

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

// export

func GroupsExport(writer io.Writer, groups []ced.Group) error {
	w := csv.NewWriter(writer)

	records := make([][]string, len(groups)+1)
	records[0] = []string{"Name", "Max Attendees", "Attendees", "Has Responded"}

	for i, g := range groups {
		records[i+1] = []string{
			string(g.Name),
			strconv.FormatUint(uint64(g.MaxAttendees), 10),
			strconv.FormatUint(uint64(g.Attendees), 10),
			strconv.FormatBool(g.HasResponded),
		}
	}

	return w.WriteAll(records)
}

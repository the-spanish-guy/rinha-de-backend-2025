package helpers

import (
	"fmt"
	"strconv"
	"time"
)

func ParseFlexibleDateTime(dateStr string) (time.Time, error) {
	layouts := []string{
		// ISO 8601 formats
		time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,           // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02T15:04:05Z",     // "2025-07-21T03:10:51Z"
		"2006-01-02T15:04:05.000Z", // "2025-07-21T03:10:51.000Z"
		"2006-01-02T15:04:05.999Z", // com milissegundos
		"2006-01-02 15:04:05",      // "2025-07-21 03:10:51"
		"2006-01-02T15:04:05",      // "2025-07-21T03:10:51"
		"2006-01-02",               // "2025-07-21"
		"02/01/2006 15:04:05",      // "21/07/2025 03:10:51"
		"02/01/2006",               // "21/07/2025"
		"1704067200",               // Unix timestamp como string
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, dateStr); err == nil {
			return parsed, nil
		}
	}

	// Unix timestamp parsing
	if timestamp, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
		return time.Unix(timestamp, 0), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s' with any known format", dateStr)
}

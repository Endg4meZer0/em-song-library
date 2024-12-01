package db

import (
	"fmt"
	"strings"
)

func arrayIntoPlaceholders(text []string, startIndex int) string {
	placeholders := make([]string, len(text))
	for i := range text {
		placeholders[i] = fmt.Sprintf("$%d", i+startIndex)
	}
	return fmt.Sprintf(`'{%s}'`, strings.Join(placeholders, ","))
}

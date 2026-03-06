package epitran

import (
	"encoding/csv"
	"strings"
)

func loadDiacritics(code string) map[string]bool {
	data, err := stripFS.ReadFile("data/strip/" + code + ".csv")
	if err != nil {
		return nil
	}
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}
	diacritics := make(map[string]bool)
	for _, record := range records {
		if len(record) >= 1 {
			diacritics[record[0]] = true
		}
	}
	if len(diacritics) == 0 {
		return nil
	}
	return diacritics
}

func stripDiacritics(text string, diacritics map[string]bool) string {
	if diacritics == nil {
		return text
	}
	var b strings.Builder
	for _, r := range text {
		if !diacritics[string(r)] {
			b.WriteRune(r)
		}
	}
	return b.String()
}

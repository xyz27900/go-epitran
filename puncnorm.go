package epitran

import (
	"encoding/csv"
	"strings"
	"sync"
)

var (
	puncnormOnce sync.Once
	puncnormMap  map[string]string
)

func loadPuncNorm() map[string]string {
	puncnormOnce.Do(func() {
		puncnormMap = make(map[string]string)
		reader := csv.NewReader(strings.NewReader(string(puncnormData)))
		records, err := reader.ReadAll()
		if err != nil {
			return
		}
		for i, record := range records {
			if i == 0 {
				continue // skip header
			}
			if len(record) >= 2 {
				puncnormMap[record[0]] = record[1]
			}
		}
	})
	return puncnormMap
}

func normPunc(text string) string {
	m := loadPuncNorm()
	var b strings.Builder
	for _, r := range text {
		s := string(r)
		if repl, ok := m[s]; ok {
			b.WriteString(repl)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

package epitran

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
	"golang.org/x/text/unicode/norm"
)

type simpleEpitran struct {
	code       string
	g2p        map[string]string // grapheme (NFD) -> phoneme (NFD)
	g2pRegex   *regexp2.Regexp   // longest-match-first anchored regex
	preProc    *ruleSet
	postProc   *ruleSet
	diacritics map[string]bool
	preproc    bool
	postproc   bool
	ligatures  bool
	tones      bool
}


var toneChars = strings.NewReplacer(
	"˩", "", "˨", "", "˧", "", "˦", "", "˥", "",
)

func newSimpleEpitran(code string, preproc, postproc, ligatures, tones bool) (*simpleEpitran, error) {
	se := &simpleEpitran{
		code:      code,
		preproc:   preproc,
		postproc:  postproc,
		ligatures: ligatures,
		tones:     tones,
	}

	var err error
	se.g2p, err = loadG2PMap(code, tones)
	if err != nil {
		return nil, err
	}

	se.g2pRegex, err = buildG2PRegex(se.g2p)
	if err != nil {
		return nil, err
	}

	// Load pre-processing rules
	if data, e := preFS.ReadFile("data/pre/" + code + ".txt"); e == nil {
		se.preProc = parseRuleFile(string(data))
	}

	// Load post-processing rules
	if data, e := postFS.ReadFile("data/post/" + code + ".txt"); e == nil {
		se.postProc = parseRuleFile(string(data))
	}

	// Load diacritics to strip
	se.diacritics = loadDiacritics(code)

	return se, nil
}

func loadG2PMap(code string, tones bool) (map[string]string, error) {
	data, err := mapFS.ReadFile("data/map/" + code + ".csv")
	if err != nil {
		return nil, &DatafileError{msg: fmt.Sprintf("cannot load G2P map for %s: %v", code, err)}
	}

	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, &DatafileError{msg: fmt.Sprintf("cannot parse G2P map for %s: %v", code, err)}
	}

	if len(records) < 1 {
		return nil, &DatafileError{msg: fmt.Sprintf("empty G2P map for %s", code)}
	}

	// Validate header
	header := records[0]
	if len(header) < 2 || header[0] != "Orth" || header[1] != "Phon" {
		return nil, &DatafileError{msg: fmt.Sprintf("G2P map for %s has invalid header: %v", code, header)}
	}

	g2p := make(map[string]string)
	graphByLine := make(map[string][]int) // track for duplicate detection

	for i, record := range records[1:] {
		if len(record) < 2 {
			return nil, &DatafileError{msg: fmt.Sprintf("G2P map for %s is malformed at line %d", code, i+2)}
		}
		grapheme := norm.NFD.String(record[0])
		phoneme := norm.NFD.String(record[1])
		if !tones {
			phoneme = toneChars.Replace(phoneme)
		}

		if _, exists := g2p[grapheme]; exists {
			graphByLine[grapheme] = append(graphByLine[grapheme], i+2)
		} else {
			g2p[grapheme] = phoneme
			graphByLine[grapheme] = append(graphByLine[grapheme], i+2)
		}
	}

	// Check for non-deterministic mappings (same grapheme on multiple lines)
	var nondets []string
	for g, lines := range graphByLine {
		if len(lines) > 1 {
			lineStrs := make([]string, len(lines))
			for i, l := range lines {
				lineStrs[i] = fmt.Sprintf("%d", l)
			}
			nondets = append(nondets, fmt.Sprintf("one-to-many G2P mapping for %q on lines %s", g, strings.Join(lineStrs, ", ")))
		}
	}
	if len(nondets) > 0 {
		return nil, &MappingError{msg: fmt.Sprintf("invalid mapping for %s:\n%s", code, strings.Join(nondets, "\n"))}
	}

	return g2p, nil
}

func buildG2PRegex(g2p map[string]string) (*regexp2.Regexp, error) {
	graphemes := make([]string, 0, len(g2p))
	for g := range g2p {
		graphemes = append(graphemes, g)
	}
	// Sort longest first; ties broken alphabetically for determinism
	sort.Slice(graphemes, func(i, j int) bool {
		li, lj := len(graphemes[i]), len(graphemes[j])
		if li != lj {
			return li > lj
		}
		return graphemes[i] < graphemes[j]
	})

	// Escape each grapheme for regex
	escaped := make([]string, len(graphemes))
	for i, g := range graphemes {
		escaped[i] = regexp2.Escape(g)
	}

	pattern := "(" + strings.Join(escaped, "|") + ")"
	re, err := regexp2.Compile(pattern, regexp2.IgnoreCase|regexp2.Unicode)
	if err != nil {
		return nil, fmt.Errorf("cannot compile G2P regex for pattern: %v", err)
	}
	return re, nil
}

func (se *simpleEpitran) transliterate(text string, normpunc bool) string {
	return se.generalTrans(text, true, normpunc)
}

func (se *simpleEpitran) strictTrans(text string, normpunc bool) string {
	return se.generalTrans(text, false, normpunc)
}

func (se *simpleEpitran) generalTrans(text string, passthrough bool, normpunc bool) string {
	// BUG FIX: Python epitran only decomposes Hangul in transliterate(), not
	// strict_trans(), which means strict_trans() always returns "" for Korean.
	// This is a bug in upstream Python epitran (simple.py lines 236-243).
	// We fix it here by decomposing in generalTrans so both paths work.
	if containsHangul(text) {
		text = decomposeHangul(text)
	}

	// NFD normalize and lowercase
	text = norm.NFD.String(strings.ToLower(text))

	// Strip language-specific diacritics
	text = stripDiacritics(text, se.diacritics)

	// Pre-processing
	if se.preproc && se.preProc != nil {
		text = se.preProc.apply(text)
	}

	// Greedy left-to-right G2P matching
	var result strings.Builder
	for len(text) > 0 {
		m, _ := se.g2pRegex.FindStringMatch(text)
		if m != nil && m.Index == 0 {
			source := m.String()
			// Lookup in map using lowercase NFD form (regex is case-insensitive)
			sourceLower := strings.ToLower(source)
			if phoneme, ok := se.g2p[sourceLower]; ok {
				result.WriteString(phoneme)
			} else if phoneme, ok := se.g2p[source]; ok {
				result.WriteString(phoneme)
			} else if passthrough {
				result.WriteString(source)
			}
			// Advance past matched runes
			text = text[len(source):]
		} else {
			// No match at position 0
			if passthrough {
				// Pass through first character (rune)
				r, size := firstRune(text)
				result.WriteRune(r)
				text = text[size:]
			} else {
				// Drop unmapped character
				_, size := firstRune(text)
				text = text[size:]
			}
		}
	}

	text = result.String()

	// Post-processing
	if se.postproc && se.postProc != nil {
		text = se.postProc.apply(text)
	}

	// Ligaturize
	if se.ligatures {
		text = Ligaturize(text)
	}

	// Punctuation normalization
	if normpunc {
		text = normPunc(text)
	}

	// NFC normalize output
	return norm.NFC.String(text)
}

func firstRune(s string) (rune, int) {
	r, size := utf8.DecodeRuneInString(s)
	return r, size
}

package epitran

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dlclark/regexp2"
	"golang.org/x/text/unicode/norm"
)

type rule struct {
	re           *regexp2.Regexp
	replacement  string // for normal rules: the b in "a -> b"
	isMetathesis bool
}

type ruleSet struct {
	rules []rule
}

// regexSub emulates Python's regex.sub(callback, text) using regexp2.
// It finds all non-overlapping matches left-to-right, calls rewrite for each,
// and builds the result string with non-matched portions preserved.
func regexSub(re *regexp2.Regexp, text string, rewrite func(*regexp2.Match) string) string {
	var b strings.Builder
	lastEnd := 0 // byte offset into text

	m, _ := re.FindStringMatch(text)
	for m != nil {
		matchStart := m.Index           // rune offset from regexp2
		matchLen := m.Length            // rune length from regexp2
		// Convert rune offsets to byte offsets
		byteStart := runeIndexToByteIndex(text, matchStart)
		byteEnd := runeIndexToByteIndex(text, matchStart+matchLen)

		// Append text between last match and this one
		b.WriteString(text[lastEnd:byteStart])
		// Append the rewritten match
		b.WriteString(rewrite(m))
		lastEnd = byteEnd

		m, _ = re.FindNextMatch(m)
	}
	// Append remaining text
	b.WriteString(text[lastEnd:])
	return b.String()
}

// runeIndexToByteIndex converts a rune offset to a byte offset in s.
func runeIndexToByteIndex(s string, runeIdx int) int {
	byteIdx := 0
	for i := 0; i < runeIdx; i++ {
		_, size := utf8.DecodeRuneInString(s[byteIdx:])
		byteIdx += size
	}
	return byteIdx
}

func parseRuleFile(content string) *ruleSet {
	rs := &ruleSet{}
	symbols := map[string]string{}

	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(rawLine)
		line = norm.NFD.String(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "%") {
			continue
		}

		// Symbol definition: ::name:: = value
		if parseSymbolDef(line, symbols) {
			continue
		}

		// Rule: a -> b / X _ Y
		r := parseRule(line, symbols)
		if r != nil {
			rs.rules = append(rs.rules, *r)
		}
	}

	return rs
}

func parseSymbolDef(line string, symbols map[string]string) bool {
	// Symbol defs must start with ::name:: at the beginning of line.
	// Pattern: ::word:: = value  (Python regex: (?P<symbol>::\w+::)\s*=\s*(?P<value>.+))
	if !strings.HasPrefix(line, "::") {
		return false
	}
	end := strings.Index(line[2:], "::")
	if end < 0 {
		return false
	}
	symbolEnd := 2 + end + 2 // byte offset past closing ::
	afterSymbol := strings.TrimSpace(line[symbolEnd:])
	if !strings.HasPrefix(afterSymbol, "=") {
		return false
	}
	symbolName := line[:symbolEnd]
	value := strings.TrimSpace(afterSymbol[1:])
	symbols[symbolName] = value
	return true
}

// findContextUnderscore finds the standalone _ that separates X from Y in the context.
// The _ must be preceded and followed by whitespace (or be at start/end of string),
// matching Python's regex behavior: \s*[_]\s*
// Uses unicode.IsSpace to handle any whitespace (upstream Python epitran data files
// contained fullwidth spaces U+3000 in pol-Latn, por-Latn, pan-Guru, tur-Latn-bab;
// we fixed those but keep the robust check).
func findContextUnderscore(context string) int {
	runes := []rune(context)
	for i, r := range runes {
		if r != '_' {
			continue
		}
		prevOk := i == 0 || unicode.IsSpace(runes[i-1])
		nextOk := i+1 >= len(runes) || unicode.IsSpace(runes[i+1])
		if prevOk && nextOk {
			// Convert rune index to byte index
			return len(string(runes[:i]))
		}
	}
	return -1
}

func subSymbols(line string, symbols map[string]string) string {
	for i := 0; i < 20; i++ { // max iterations to resolve nested symbols
		found := false
		for name, value := range symbols {
			if strings.Contains(line, name) {
				line = strings.ReplaceAll(line, name, value)
				found = true
			}
		}
		if !found {
			break
		}
	}
	return line
}

func parseRule(line string, symbols map[string]string) *rule {
	// Parse: a -> b / X _ Y
	arrowIdx := strings.Index(line, "->")
	if arrowIdx < 0 {
		return nil
	}

	// Find the "/" separator after the arrow
	rest := line[arrowIdx+2:]
	slashIdx := strings.Index(rest, "/")
	if slashIdx < 0 {
		return nil
	}
	slashIdx += arrowIdx + 2

	a := strings.TrimSpace(line[:arrowIdx])
	b := strings.TrimSpace(line[arrowIdx+2 : slashIdx])
	context := strings.TrimSpace(line[slashIdx+1:])

	// Substitute symbols in the full context BEFORE splitting on _,
	// since symbol names like ::pal_consonant:: contain underscores.
	a = subSymbols(a, symbols)
	b = subSymbols(b, symbols)
	context = subSymbols(context, symbols)

	// Find the standalone _ separator (surrounded by whitespace or at edges)
	underIdx := findContextUnderscore(context)
	if underIdx < 0 {
		return nil
	}

	X := strings.TrimSpace(context[:underIdx])
	Y := strings.TrimSpace(context[underIdx+1:])

	// Replace # with ^ (start) or $ (end)
	X = strings.ReplaceAll(X, "#", "^")
	Y = strings.ReplaceAll(Y, "#", "$")

	// Replace 0 with empty string (0 means null/empty in phonological rules)
	a = strings.ReplaceAll(a, "0", "")
	b = strings.ReplaceAll(b, "0", "")

	// Check for metathesis (sw1/sw2 groups in a)
	isMetathesis := strings.Contains(a, "(?P<sw1>") && strings.Contains(a, "(?P<sw2>")

	// Convert Python named groups (?P<name>...) to .NET syntax (?<name>...) for regexp2
	a = pythonToNetNamedGroups(a)

	// Build the full regex pattern using .NET named group syntax
	var pattern string
	if isMetathesis {
		pattern = fmt.Sprintf("(?<X>%s)%s(?<Y>%s)", X, a, Y)
	} else {
		pattern = fmt.Sprintf("(?<X>%s)(?<a>%s)(?<Y>%s)", X, a, Y)
	}

	// Fix Unicode property names for regexp2 compatibility
	pattern = fixUnicodeProperties(pattern)

	re, err := regexp2.Compile(pattern, regexp2.Unicode)
	if err != nil {
		return nil // skip rules that can't compile
	}

	return &rule{
		re:           re,
		replacement:  b,
		isMetathesis: isMetathesis,
	}
}

func (rs *ruleSet) apply(text string) string {
	for _, r := range rs.rules {
		text = r.apply(text)
	}
	return text
}

func (r *rule) apply(text string) string {
	return regexSub(r.re, text, func(m *regexp2.Match) string {
		if r.isMetathesis {
			X := groupByName(m, "X")
			sw1 := groupByName(m, "sw1")
			sw2 := groupByName(m, "sw2")
			Y := groupByName(m, "Y")
			return X + sw2 + sw1 + Y
		}
		X := groupByName(m, "X")
		Y := groupByName(m, "Y")
		return X + r.replacement + Y
	})
}

// pythonToNetNamedGroups converts Python-style (?P<name>...) to .NET-style (?<name>...).
func pythonToNetNamedGroups(s string) string {
	return strings.ReplaceAll(s, "(?P<", "(?<")
}

// fixUnicodeProperties converts long Unicode property names to short forms
// supported by regexp2 (which uses .NET regex syntax).
func fixUnicodeProperties(s string) string {
	return strings.ReplaceAll(s, `\p{Letter}`, `\p{L}`)
}

func groupByName(m *regexp2.Match, name string) string {
	g := m.GroupByName(name)
	if g == nil {
		return ""
	}
	return g.String()
}

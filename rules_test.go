package epitran

import (
	"testing"

	"github.com/dlclark/regexp2"
	"golang.org/x/text/unicode/norm"
)

func TestRegexSub(t *testing.T) {
	// Test basic regex sub functionality (.NET named group syntax)
	re, err := regexp2.Compile(`(?<X>)(?<a>j)(?<Y>)`, regexp2.Unicode)
	if err != nil {
		t.Fatal(err)
	}
	result := regexSub(re, "joukov", func(m *regexp2.Match) string {
		X := groupByName(m, "X")
		Y := groupByName(m, "Y")
		return X + "ʒ" + Y
	})
	if result != "ʒoukov" {
		t.Errorf("regexSub j->ʒ: got %q, want %q", result, "ʒoukov")
	}
}

func TestFrenchPreproc(t *testing.T) {
	// Load French pre-processing rules
	data, err := preFS.ReadFile("data/pre/fra-Latn.txt")
	if err != nil {
		t.Fatal(err)
	}
	rs := parseRuleFile(string(data))
	t.Logf("Loaded %d French pre-processing rules", len(rs.rules))

	tests := []struct {
		input, want string
	}{
		{"joukov", "ʒukɔv"},
		{"paris", "paʀi"},
	}

	for _, tt := range tests {
		got := rs.apply(tt.input)
		t.Logf("preproc(%q) = %q (want %q)", tt.input, got, tt.want)
	}
}

func TestRulePattern(t *testing.T) {
	// Test a simple rule: j -> ʒ / _
	// This means: replace j with ʒ in any context
	// Pattern: (?P<X>)(?P<a>j)(?P<Y>)
	symbols := map[string]string{}
	r := parseRule("j -> ʒ / _", symbols)
	if r == nil {
		t.Fatal("parseRule returned nil")
	}
	t.Logf("Rule pattern: %s", r.re.String())

	result := r.apply("joukov")
	t.Logf("j -> ʒ applied to 'joukov': %q", result)
	if result != "ʒoukov" {
		t.Errorf("got %q, want %q", result, "ʒoukov")
	}
}

func TestRuleWithContext(t *testing.T) {
	// Test: s -> z / (a|e|i|o|u) _ (a|e|i|o|u)
	// intervocalic voicing
	symbols := map[string]string{
		"::vowel::": "a|e|i|o|u",
	}
	r := parseRule("s -> z / (::vowel::) _ (::vowel::)", symbols)
	if r == nil {
		t.Fatal("parseRule returned nil")
	}
	t.Logf("Rule pattern: %s", r.re.String())

	result := r.apply("casa")
	t.Logf("intervocalic s -> z in 'casa': %q", result)
	if result != "caza" {
		t.Errorf("got %q, want %q", result, "caza")
	}
}

func TestRuleWithHash(t *testing.T) {
	// Test: s -> 0 / (a|e|i|o|u).+ _ #
	// delete final s after vowel
	symbols := map[string]string{
		"::vowel::": "a|e|i|o|u",
	}
	r := parseRule("s -> 0 / (::vowel::).+ _ #", symbols)
	if r == nil {
		t.Fatal("parseRule returned nil")
	}
	t.Logf("Rule pattern: %s", r.re.String())

	result := r.apply("paris")
	t.Logf("final s deletion in 'paris': %q", result)
	if result != "pari" {
		t.Errorf("got %q, want %q", result, "pari")
	}
}

func TestHindiPostproc(t *testing.T) {
	data, err := postFS.ReadFile("data/post/hin-Deva.txt")
	if err != nil {
		t.Fatal(err)
	}
	rs := parseRuleFile(string(data))
	t.Logf("Loaded %d Hindi post-processing rules", len(rs.rules))
	for i, r := range rs.rules {
		pat := r.re.String()
		if len(pat) > 120 {
			pat = pat[:120] + "..."
		}
		t.Logf("  rule[%d] metathesis=%v repl=%q pattern=%s", i, r.isMetathesis, r.replacement, pat)
	}

	// Trace each rule on the Hindi text
	text := "nəməsə\u094Dtəe" // nəməsə्təe
	t.Logf("Input: %q", text)
	for i, r := range rs.rules {
		newText := r.apply(text)
		if newText != text {
			t.Logf("  Rule %d: %q -> %q (pattern: %s)", i, text, newText, r.re.String()[:min(100, len(r.re.String()))])
		}
		text = newText
	}
	t.Logf("Final: %q", text)

	want := "nəmste"
	if text != want {
		t.Errorf("Hindi postproc: got %q, want %q", text, want)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestPolishPreproc(t *testing.T) {
	data, err := preFS.ReadFile("data/pre/pol-Latn.txt")
	if err != nil {
		t.Fatal(err)
	}
	rs := parseRuleFile(string(data))
	t.Logf("Loaded %d Polish pre-processing rules", len(rs.rules))

	// Apply to NFD-normalized lowercase text (must be NFD for rules to work)
	text := norm.NFD.String("dziękuję")
	t.Logf("Input: %q", text)
	for i, r := range rs.rules {
		newText := r.apply(text)
		if newText != text {
			t.Logf("  Rule %d: %q -> %q", i, text, newText)
		}
		text = newText
	}
	t.Logf("Final: %q", text)

	want := "d͡ʑɛŋkujɛ"
	if text != want {
		t.Errorf("Polish preproc: got %q, want %q", text, want)
	}
}

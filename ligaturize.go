package epitran

import "strings"

var ligatureMap = []struct{ from, to string }{
	{"t͡s", "ʦ"},
	{"t͡ʃ", "ʧ"},
	{"t͡ɕ", "ʨ"},
	{"d͡z", "ʣ"},
	{"d͡ʒ", "ʤ"},
	{"d͡ʑ", "ʥ"},
}

// Ligaturize converts standard IPA affricate sequences to precomposed ligatures.
func Ligaturize(text string) string {
	for _, m := range ligatureMap {
		text = strings.ReplaceAll(text, m.from, m.to)
	}
	return text
}

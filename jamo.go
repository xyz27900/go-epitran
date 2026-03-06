package epitran

import "strings"

// Hangul syllable decomposition to Compatibility Jamo.
// Korean epitran maps expect Compatibility Jamo (U+3131-U+3163),
// NOT Conjoining Jamo (U+1100-U+11FF) which NFD produces.

var leadToCompat = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ',
	'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

var vowelToCompat = []rune{
	'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ',
	'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ',
}

var tailToCompat = []rune{
	0, // index 0 = no tail
	'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ',
	'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ',
	'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

func decomposeHangul(text string) string {
	var b strings.Builder
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7A3 {
			offset := r - 0xAC00
			lead := offset / (21 * 28)
			vowel := (offset % (21 * 28)) / 28
			tail := offset % 28

			b.WriteRune(leadToCompat[lead])
			b.WriteRune(vowelToCompat[vowel])
			if tail > 0 {
				b.WriteRune(tailToCompat[tail])
			}
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func containsHangul(text string) bool {
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7A3 {
			return true
		}
	}
	return false
}

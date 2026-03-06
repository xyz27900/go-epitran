//go:build python

package epitran

import (
	"os/exec"
	"strings"
	"testing"
)

func pythonTransliterate(t *testing.T, code, word string) string {
	t.Helper()
	script := "import epitran; print(epitran.Epitran('" + code + "').transliterate('" + word + "'), end='')"
	out, err := exec.Command("python3", "-c", script).Output()
	if err != nil {
		t.Fatalf("python3 epitran(%q, %q): %v", code, word, err)
	}
	return string(out)
}

func TestTransliteratePython(t *testing.T) {
	tests := []struct {
		code string
		word string
	}{
		// Russian
		{"rus-Cyrl", "привет"},
		{"rus-Cyrl", "спасибо"},
		{"rus-Cyrl", "россия"},
		{"rus-Cyrl", "красный"},
		{"rus-Cyrl", "здравствуйте"},
		{"rus-Cyrl", "москва"},
		{"rus-Cyrl", "сергей"},
		{"rus-Cyrl", "путин"},

		// Ukrainian
		{"ukr-Cyrl", "привіт"},
		{"ukr-Cyrl", "дякую"},
		{"ukr-Cyrl", "водограй"},
		{"ukr-Cyrl", "київ"},
		{"ukr-Cyrl", "зеленський"},
		{"ukr-Cyrl", "україна"},
		{"ukr-Cyrl", "шевченко"},

		// German
		{"deu-Latn", "kindergarten"},
		{"deu-Latn", "volkswagen"},
		{"deu-Latn", "gesundheit"},
		{"deu-Latn", "übung"},
		{"deu-Latn", "zwölf"},
		{"deu-Latn", "pflicht"},
		{"deu-Latn", "natürlich"},
		{"deu-Latn", "wissenschaft"},
		{"deu-Latn", "berlin"},
		{"deu-Latn", "münchen"},
		{"deu-Latn", "straße"},
		{"deu-Latn", "schön"},
		{"deu-Latn", "sprechen"},
		{"deu-Latn", "ich"},
		{"deu-Latn", "ach"},
		{"deu-Latn", "zeit"},
		{"deu-Latn", "pferd"},

		// French
		{"fra-Latn", "merci"},
		{"fra-Latn", "beaucoup"},
		{"fra-Latn", "fromage"},
		{"fra-Latn", "fenêtre"},
		{"fra-Latn", "oiseau"},
		{"fra-Latn", "feuille"},
		{"fra-Latn", "champagne"},
		{"fra-Latn", "paris"},
		{"fra-Latn", "bonjour"},
		{"fra-Latn", "joukov"},
		{"fra-Latn", "château"},
		{"fra-Latn", "français"},
		{"fra-Latn", "eau"},
		{"fra-Latn", "chien"},
		{"fra-Latn", "fille"},

		// Polish
		{"pol-Latn", "dziękuję"},
		{"pol-Latn", "przepraszam"},
		{"pol-Latn", "cześć"},
		{"pol-Latn", "źródło"},
		{"pol-Latn", "chrząszcz"},
		{"pol-Latn", "warszawa"},
		{"pol-Latn", "szczecin"},
		{"pol-Latn", "łódź"},

		// Portuguese
		{"por-Latn", "obrigado"},
		{"por-Latn", "portugal"},
		{"por-Latn", "são"},
		{"por-Latn", "manhã"},
		{"por-Latn", "lisboa"},
		{"por-Latn", "brasil"},
		{"por-Latn", "coração"},

		// Hindi
		{"hin-Deva", "नमस्कार"},
		{"hin-Deva", "प्रणाम"},
		{"hin-Deva", "धन्यवाद"},
		{"hin-Deva", "अच्छा"},
		{"hin-Deva", "सत्यम्"},
		{"hin-Deva", "राहुल"},
		{"hin-Deva", "नमस्ते"},
		{"hin-Deva", "भारत"},
		{"hin-Deva", "दिल्ली"},

		// Thai
		{"tha-Thai", "ประเทศไทย"},
		{"tha-Thai", "ภาษา"},
		{"tha-Thai", "มาลัย"},
		{"tha-Thai", "กรุงเทพ"},
		{"tha-Thai", "สวัสดี"},
		{"tha-Thai", "ขอบคุณ"},

		// Korean
		{"kor-Hang", "대한민국"},
		{"kor-Hang", "안녕하세요"},
		{"kor-Hang", "감사합니다"},
		{"kor-Hang", "사랑"},
		{"kor-Hang", "김민수"},
		{"kor-Hang", "서울"},
		{"kor-Hang", "한국"},

		// Spanish
		{"spa-Latn", "gracias"},
		{"spa-Latn", "pronunciación"},
		{"spa-Latn", "guerrilla"},
		{"spa-Latn", "ahora"},
		{"spa-Latn", "lluvia"},
		{"spa-Latn", "españa"},
		{"spa-Latn", "barcelona"},
		{"spa-Latn", "corazón"},

		// Italian
		{"ita-Latn", "grazie"},
		{"ita-Latn", "ciao"},
		{"ita-Latn", "buonasera"},
		{"ita-Latn", "famiglia"},
		{"ita-Latn", "spaghetti"},
		{"ita-Latn", "roma"},
		{"ita-Latn", "firenze"},
		{"ita-Latn", "buongiorno"},

		// Swedish
		{"swe-Latn", "tack"},
		{"swe-Latn", "skön"},
		{"swe-Latn", "björk"},

		// Czech
		{"ces-Latn", "děkuji"},
		{"ces-Latn", "křižovatka"},
		{"ces-Latn", "příliš"},
		{"ces-Latn", "praha"},

		// Hungarian
		{"hun-Latn", "viszontlátásra"},
		{"hun-Latn", "egészségedre"},
		{"hun-Latn", "budapest"},
		{"hun-Latn", "magyar"},
		{"hun-Latn", "köszönöm"},

		// Dutch
		{"nld-Latn", "dank"},
		{"nld-Latn", "uitstekend"},
		{"nld-Latn", "schiphol"},
		{"nld-Latn", "amsterdam"},

		// Arabic
		{"ara-Arab", "شكرا"},
		{"ara-Arab", "مصر"},
		{"ara-Arab", "لبنان"},
		{"ara-Arab", "محمد"},
		{"ara-Arab", "بغداد"},

		// Persian
		{"fas-Arab", "سلام"},
		{"fas-Arab", "خداحافظ"},
		{"fas-Arab", "مرسی"},
		{"fas-Arab", "تهران"},
		{"fas-Arab", "ایران"},
		{"fas-Arab", "فارسی"},

		// Bengali
		{"ben-Beng", "নমস্কার"},
		{"ben-Beng", "ধন্যবাদ"},
		{"ben-Beng", "বাংলাদেশ"},
		{"ben-Beng", "বাংলা"},
		{"ben-Beng", "কলকাতা"},
		{"ben-Beng", "ঢাকা"},

		// Amharic
		{"amh-Ethi", "ሰላም"},
		{"amh-Ethi", "አመሰግናለሁ"},
		{"amh-Ethi", "አዲስ"},

		// Vietnamese
		{"vie-Latn", "xin"},
		{"vie-Latn", "chào"},
		{"vie-Latn", "cảm"},
		{"vie-Latn", "ơn"},
		{"vie-Latn", "hanoi"},

		// Turkish
		{"tur-Latn", "merhaba"},
		{"tur-Latn", "teşekkürler"},
		{"tur-Latn", "hoşçakalın"},
		{"tur-Latn", "istanbul"},
		{"tur-Latn", "ankara"},
		{"tur-Latn", "erdoğan"},
		{"tur-Latn", "türkiye"},
		{"tur-Latn", "güzel"},
		{"tur-Latn", "çiçek"},
		{"tur-Latn", "şehir"},

		// Urdu
		{"urd-Arab", "محمد"},
		{"urd-Arab", "اقبال"},
		{"urd-Arab", "پاکستان"},

		// Kazakh
		{"kaz-Cyrl", "астана"},
		{"kaz-Cyrl", "қазақстан"},
	}

	// Deduplicate by code to create one Epitran instance per language
	instances := map[string]*Epitran{}

	for _, tt := range tests {
		t.Run(tt.code+"/"+tt.word, func(t *testing.T) {
			epi, ok := instances[tt.code]
			if !ok {
				var err error
				epi, err = New(tt.code)
				if err != nil {
					t.Fatalf("New(%q): %v", tt.code, err)
				}
				instances[tt.code] = epi
			}

			got := epi.Transliterate(tt.word)
			want := pythonTransliterate(t, tt.code, tt.word)

			if got != want {
				// Show hex for debugging Unicode differences
				gotHex := strings.Map(func(r rune) rune { return r }, got)
				wantHex := strings.Map(func(r rune) rune { return r }, want)
				t.Errorf("Transliterate(%q, %q):\n  got  = %q (%x)\n  want = %q (%x)", tt.code, tt.word, gotHex, []byte(got), wantHex, []byte(want))
			}
		})
	}
}

package epitran

import "testing"

func TestTransliterate(t *testing.T) {
	tests := []struct {
		code  string
		input string
		want  string
	}{
		// Russian
		{"rus-Cyrl", "Москва", "moskva"},
		{"rus-Cyrl", "сергей", "sʲerɡej"},
		{"rus-Cyrl", "путин", "putʲin"},

		// Ukrainian
		{"ukr-Cyrl", "київ", "kɪjiv"},
		{"ukr-Cyrl", "зеленський", "zɛlɛnsʲkɪj"},
		{"ukr-Cyrl", "україна", "ukrɑjinɑ"},
		{"ukr-Cyrl", "шевченко", "ʂɛvʈ͡ʂɛnkɔ"},

		// Arabic
		{"ara-Arab", "محمد", "mħmd"},
		{"ara-Arab", "بغداد", "bɣdaːd"},

		// Turkish
		{"tur-Latn", "istanbul", "istanbul"},
		{"tur-Latn", "ankara", "ankaɾa"},
		{"tur-Latn", "erdoğan", "eɾdoɰan"},
		{"tur-Latn", "türkiye", "tyɾkije"},
		{"tur-Latn", "güzel", "ɡyzel"},
		{"tur-Latn", "çiçek", "t͡ʃit͡ʃek"},
		{"tur-Latn", "şehir", "ʃehiɾ"},

		// French
		{"fra-Latn", "paris", "paʀi"},
		{"fra-Latn", "bonjour", "bɔ̃ʒuʀ"},
		{"fra-Latn", "joukov", "ʒukɔv"},
		{"fra-Latn", "château", "ʃato"},
		{"fra-Latn", "français", "fʀɑ̃sɛs"},
		{"fra-Latn", "eau", "o"},
		{"fra-Latn", "chien", "ʃjɛ̃"},
		{"fra-Latn", "fille", "fij"},

		// German
		{"deu-Latn", "berlin", "bərliːn"},
		{"deu-Latn", "münchen", "mynxn̩"},
		{"deu-Latn", "straße", "ʃtraːsə"},
		{"deu-Latn", "schön", "ʃøːn"},
		{"deu-Latn", "sprechen", "ʃprɛçn̩"},
		{"deu-Latn", "ich", "ɪç"},
		{"deu-Latn", "ach", "ax"},
		{"deu-Latn", "zeit", "t͡saɪ̯t"},
		{"deu-Latn", "pferd", "p͡fɛːɐ̯t"},

		// Polish
		{"pol-Latn", "warszawa", "varʂava"},
		{"pol-Latn", "szczecin", "ʂt͡ʂɛt͡ɕin"},
		{"pol-Latn", "łódź", "wut͡ɕ"},

		// Portuguese
		{"por-Latn", "lisboa", "liʃbowɐ"},
		{"por-Latn", "brasil", "bɾɐzil"},
		{"por-Latn", "coração", "koɾɐsɐ̃w̃"},

		// Hindi
		{"hin-Deva", "राहुल", "raːɦul"},
		{"hin-Deva", "नमस्ते", "nəmste"},
		{"hin-Deva", "भारत", "b̤aːrət"},
		{"hin-Deva", "दिल्ली", "dilliː"},

		// Bengali
		{"ben-Beng", "বাংলা", "baŋl̪a"},
		{"ben-Beng", "কলকাতা", "kɔl̪ɔkat̪a"},
		{"ben-Beng", "ঢাকা", "d̤aka"},

		// Thai
		{"tha-Thai", "กรุงเทพ", "kruŋtʰeːp"},
		{"tha-Thai", "สวัสดี", "sawatdiː"},
		{"tha-Thai", "ขอบคุณ", "kʰɔːbkʰun"},

		// Hungarian
		{"hun-Latn", "budapest", "budɒpɛʃt"},
		{"hun-Latn", "magyar", "mɒɟɒr"},
		{"hun-Latn", "köszönöm", "køsønøm"},

		// Persian
		{"fas-Arab", "تهران", "thrɒn"},
		{"fas-Arab", "ایران", "ɒjrɒn"},
		{"fas-Arab", "فارسی", "fɒrsj"},

		// Korean
		{"kor-Hang", "김민수", "kimminsu"},
		{"kor-Hang", "서울", "sʌul"},
		{"kor-Hang", "한국", "hankuk"},

		// Spanish
		{"spa-Latn", "españa", "espaɲa"},
		{"spa-Latn", "barcelona", "baɾselona"},
		{"spa-Latn", "corazón", "koɾason"},

		// Italian
		{"ita-Latn", "roma", "roma"},
		{"ita-Latn", "firenze", "firense"},
		{"ita-Latn", "buongiorno", "buond͡ʒorno"},

		// Czech
		{"ces-Latn", "praha", "braɦa"},

		// Dutch
		{"nld-Latn", "amsterdam", "ɑmstɛrdɑm"},

		// Vietnamese
		{"vie-Latn", "hanoi", "hanɔj"},

		// Amharic
		{"amh-Ethi", "አዲስ", "ədis"},
	}

	for _, tt := range tests {
		t.Run(tt.code+"/"+tt.input, func(t *testing.T) {
			epi, err := New(tt.code)
			if err != nil {
				t.Fatalf("New(%q): %v", tt.code, err)
			}
			got := epi.Transliterate(tt.input)
			if got != tt.want {
				t.Errorf("Transliterate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStrictTrans(t *testing.T) {
	tests := []struct {
		code  string
		input string
		want  string
	}{
		{"rus-Cyrl", "Москва", "moskva"},
		{"fra-Latn", "bonjour", "bɔ̃ʒuʀ"},
		{"deu-Latn", "münchen", "mynxn̩"},
		{"hin-Deva", "राहुल", "raːɦul"},
		// BUG FIX: Python epitran strict_trans returns "" for Korean because it
		// doesn't decompose Hangul. We fix this — strict_trans should work like
		// transliterate but dropping unmapped chars.
		{"kor-Hang", "김민수", "kimminsu"},
		{"kor-Hang", "서울", "sʌul"},
	}

	for _, tt := range tests {
		t.Run(tt.code+"/"+tt.input, func(t *testing.T) {
			epi, err := New(tt.code)
			if err != nil {
				t.Fatalf("New(%q): %v", tt.code, err)
			}
			got := epi.StrictTrans(tt.input)
			if got != tt.want {
				t.Errorf("StrictTrans(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSpecialLanguagesReturnError(t *testing.T) {
	for _, code := range []string{"eng-Latn", "cmn-Hans", "cmn-Hant", "jpn-Jpan", "yue-Hant"} {
		_, err := New(code)
		if err != ErrUnsupportedSpecialLanguage {
			t.Errorf("New(%q) err = %v, want ErrUnsupportedSpecialLanguage", code, err)
		}
	}
}

func TestLigatures(t *testing.T) {
	epi, err := New("tur-Latn", WithLigatures(true))
	if err != nil {
		t.Fatal(err)
	}
	got := epi.Transliterate("çiçek")
	want := "ʧiʧek"
	if got != want {
		t.Errorf("Transliterate with ligatures = %q, want %q", got, want)
	}
}

func TestWithoutPreproc(t *testing.T) {
	epi, err := New("fra-Latn", WithPreproc(false))
	if err != nil {
		t.Fatal(err)
	}
	got := epi.Transliterate("paris")
	// Without preprocessing, the French rules that convert final s -> 0 won't apply
	// so we expect different output
	t.Logf("fra-Latn 'paris' without preproc = %q", got)
}

func TestWithoutPostproc(t *testing.T) {
	epi, err := New("hin-Deva", WithPostproc(false))
	if err != nil {
		t.Fatal(err)
	}
	got := epi.Transliterate("राहुल")
	// Without postprocessing, the schwa-deletion rules won't apply
	t.Logf("hin-Deva 'राहुल' without postproc = %q", got)
}

func TestHangulDecomposition(t *testing.T) {
	got := decomposeHangul("가")
	want := "ㄱㅏ"
	if got != want {
		t.Errorf("decomposeHangul(가) = %q, want %q", got, want)
	}
}

func TestLigaturizeFunc(t *testing.T) {
	tests := []struct{ in, want string }{
		{"t͡s", "ʦ"},
		{"t͡ʃ", "ʧ"},
		{"d͡ʒ", "ʤ"},
		{"abc", "abc"},
	}
	for _, tt := range tests {
		got := Ligaturize(tt.in)
		if got != tt.want {
			t.Errorf("Ligaturize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

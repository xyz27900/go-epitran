package epitran

import "testing"

func TestDifferentialParity(t *testing.T) {
	// All test vectors generated from Python epitran for byte-identical parity
	tests := []struct {
		code  string
		input string
		want  string
	}{
		// Russian
		{"rus-Cyrl", "привет", "prʲivʲet"},
		{"rus-Cyrl", "спасибо", "spasʲibo"},
		{"rus-Cyrl", "россия", "rossʲia"},
		{"rus-Cyrl", "красный", "krasnɨj"},
		{"rus-Cyrl", "здравствуйте", "zdravstvujtʲe"},

		// Ukrainian
		{"ukr-Cyrl", "привіт", "prɪvit"},
		{"ukr-Cyrl", "дякую", "dʲɑkuu"},
		{"ukr-Cyrl", "водограй", "vɔdɔɦrɑj"},

		// German
		{"deu-Latn", "kindergarten", "kɪndɛːrɡaːrtən"},
		{"deu-Latn", "volkswagen", "fɔlksvaɡn̩"},
		{"deu-Latn", "gesundheit", "ɡəzʊndhaɪ̯t"},
		{"deu-Latn", "übung", "yːbʊŋk"},
		{"deu-Latn", "zwölf", "t͡svœlf"},
		{"deu-Latn", "pflicht", "p͡flɪçt"},
		{"deu-Latn", "natürlich", "naːtyːrlɪç"},
		{"deu-Latn", "wissenschaft", "vɪzɛnʃaft"},

		// French
		{"fra-Latn", "merci", "mɛʀsi"},
		{"fra-Latn", "beaucoup", "bəoku"},
		{"fra-Latn", "fromage", "fʀɔmaʒ"},
		{"fra-Latn", "fenêtre", "fənə̂tʀ"},
		{"fra-Latn", "oiseau", "wazo"},
		{"fra-Latn", "feuille", "fœjl"},
		{"fra-Latn", "champagne", "ʃɑ̃paɲ"},

		// Polish
		{"pol-Latn", "dziękuję", "d͡ʑɛŋkujɛ"},
		{"pol-Latn", "przepraszam", "pʂɛpraʂam"},
		{"pol-Latn", "cześć", "t͡ʂɛɕt͡ɕ"},
		{"pol-Latn", "źródło", "ʑrudwɔ"},
		{"pol-Latn", "chrząszcz", "xʂɔ̃ʂt͡ʂ"},

		// Portuguese
		{"por-Latn", "obrigado", "obʁiɡɐdo"},
		{"por-Latn", "portugal", "poɾtuɡɐl"},
		{"por-Latn", "são", "sɐ̃w̃"},
		{"por-Latn", "manhã", "mɐnɐ̃"},

		// Hindi
		{"hin-Deva", "नमस्कार", "nəmskaːr"},
		{"hin-Deva", "प्रणाम", "prəɳaːm"},
		{"hin-Deva", "धन्यवाद", "d̤ənjəvaːd"},
		{"hin-Deva", "अच्छा", "at͡ʃt͡ʃʰaː"},
		{"hin-Deva", "सत्यम्", "sətjəm"},

		// Thai
		{"tha-Thai", "ประเทศไทย", "pratʰeːttʰaj"},
		{"tha-Thai", "ภาษา", "pʰaːsaː"},
		{"tha-Thai", "มาลัย", "maːlaj"},

		// Korean
		{"kor-Hang", "대한민국", "tɛanminkuk"},
		{"kor-Hang", "안녕하세요", "annjʌŋhasejo"},
		{"kor-Hang", "감사합니다", "kamsaamnita"},
		{"kor-Hang", "사랑", "salaŋ"},

		// Spanish
		{"spa-Latn", "gracias", "ɡɾasias"},
		{"spa-Latn", "pronunciación", "pɾonunsiasion"},
		{"spa-Latn", "guerrilla", "ɡeriʝa"},
		{"spa-Latn", "ahora", "aoɾa"},
		{"spa-Latn", "lluvia", "ʝubja"},

		// Italian
		{"ita-Latn", "grazie", "ɡrasie"},
		{"ita-Latn", "ciao", "t͡ʃao"},
		{"ita-Latn", "buonasera", "buonasera"},
		{"ita-Latn", "famiglia", "famiʎa"},
		{"ita-Latn", "spaghetti", "spaɡetːi"},

		// Swedish
		{"swe-Latn", "tack", "tɐk"},
		{"swe-Latn", "skön", "sɕøːn"},
		{"swe-Latn", "björk", "bjœrk"},

		// Czech
		{"ces-Latn", "děkuji", "djekujɪ"},
		{"ces-Latn", "křižovatka", "ɡrɪʒovatka"},
		{"ces-Latn", "příliš", "briːlɪʃ"},

		// Hungarian
		{"hun-Latn", "viszontlátásra", "visontlaːtaːʃrɒ"},
		{"hun-Latn", "egészségedre", "ɛɡeːsʃeːɡɛdrɛ"},

		// Dutch
		{"nld-Latn", "dank", "dɑnk"},
		{"nld-Latn", "uitstekend", "œɥtsteːkɛnt"},
		{"nld-Latn", "schiphol", "sxɪfɔl"},

		// Arabic
		{"ara-Arab", "شكرا", "ʃkraː"},
		{"ara-Arab", "مصر", "msˤr"},
		{"ara-Arab", "لبنان", "lbnaːn"},

		// Persian
		{"fas-Arab", "سلام", "slɒm"},
		{"fas-Arab", "خداحافظ", "xdɒhɒfz"},
		{"fas-Arab", "مرسی", "mrsj"},

		// Bengali
		{"ben-Beng", "নমস্কার", "n̪ɔmɔs̪kar"},
		{"ben-Beng", "ধন্যবাদ", "d̪̤ɔn̪d͡zɔbad̪"},
		{"ben-Beng", "বাংলাদেশ", "baŋl̪ad̪eɕ"},

		// Amharic
		{"amh-Ethi", "ሰላም", "səlam"},
		{"amh-Ethi", "አመሰግናለሁ", "əməsəɡɨnaləhu"},

		// Vietnamese
		{"vie-Latn", "xin", "sin"},
		{"vie-Latn", "chào", "caw"},
		{"vie-Latn", "cảm", "kam"},
		{"vie-Latn", "ơn", "ʔɤn"},

		// Turkish
		{"tur-Latn", "merhaba", "meɾhaba"},
		{"tur-Latn", "teşekkürler", "teʃekkyɾleɾ"},
		{"tur-Latn", "hoşçakalın", "hoʃt͡ʃakalɯn"},
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

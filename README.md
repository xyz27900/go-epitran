# go-epitran

A Go port of the Python [epitran](https://github.com/dmort27/epitran) library — grapheme-to-phoneme (G2P) transliteration for 153+ languages.

All 302 data files are embedded in the binary via `//go:embed`, so there are no external dependencies at runtime.

## Install

```
go get github.com/xyz27900/go-epitran
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/xyz27900/go-epitran"
)

func main() {
	epi, err := epitran.New("rus-Cyrl")
	if err != nil {
		panic(err)
	}

	fmt.Println(epi.Transliterate("Москва"))    // moskva
	fmt.Println(epi.Transliterate("сергей"))     // sʲerɡej
}
```

## API

| Function                                             | Description                                         |
|------------------------------------------------------|-----------------------------------------------------|
| `New(code string, opts ...Option) (*Epitran, error)` | Create an instance for a language-script code       |
| `Transliterate(text string) string`                  | Convert text to IPA (unmapped chars passed through) |
| `TransliterateNorm(text string) string`              | Same as above with punctuation normalization        |
| `StrictTrans(text string) string`                    | Convert to IPA, dropping unmapped chars             |
| `StrictTransNorm(text string) string`                | Strict mode with punctuation normalization          |

### Options

```go
epitran.New("fra-Latn",
	epitran.WithPreproc(false),    // disable preprocessing rules
	epitran.WithPostproc(false),   // disable postprocessing rules
	epitran.WithLigatures(true),   // output affricate ligatures (e.g. t͡ʃ → ʧ)
	epitran.WithTones(true),       // include tone marks
)
```

## Supported languages

The full list of supported language-script codes (e.g. `rus-Cyrl`, `ara-Arab`, `hin-Deva`) matches [upstream epitran](https://github.com/dmort27/epitran#language-support).

Five special codes require handlers not yet ported and return `ErrUnsupportedSpecialLanguage`: `eng-Latn`, `cmn-Hans`, `cmn-Hant`, `jpn-Jpan`, `yue-Hant`.

## Python parity

All output is verified byte-for-byte against Python epitran across 21 languages:

```
pip install epitran
go test -v -tags python ./...
```

## License

MIT — see [LICENSE](LICENSE).

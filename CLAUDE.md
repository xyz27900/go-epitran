# go-epitran

Go port of Python [epitran](https://github.com/dmort27/epitran) — grapheme-to-phoneme transliteration for 153+ languages.

## Architecture

Pipeline: NFD → lowercase → strip diacritics → pre-rules → G2P lookup → post-rules → ligaturize → puncnorm → NFC

## Key dependency

`github.com/dlclark/regexp2` for .NET-style regex (lookaheads, named groups) since Go stdlib `regexp` lacks these.

## Data

302 embedded files (maps, pre/post rules, strip lists, puncnorm) via `//go:embed`.

## Critical gotchas

- `(?P<name>...)` → `(?<name>...)` conversion for regexp2
- Symbol defs (`::name::`) must start at line beginning (not match `(?=...)`)
- Underscore separator uses `unicode.IsSpace` (fullwidth spaces in data files)
- `\p{Letter}` → `\p{L}` for regexp2
- Korean Hangul decomposition happens in `generalTrans` (bug fix vs Python which only does it in `transliterate`)

## 5 unsupported special language codes

`eng-Latn`, `cmn-Hans`, `cmn-Hant`, `jpn-Jpan`, `yue-Hant` — these require special handlers not yet ported.

## Python parity tests

```
go test -v -tags python ./...
```

Requires `pip install epitran`.

// Package epitran provides grapheme-to-phoneme transliteration for 153+ languages.
// It is a Go port of the Python epitran library (https://github.com/dmort27/epitran).
package epitran

import "embed"

//go:embed data/map/*.csv
var mapFS embed.FS

//go:embed data/pre/*.txt
var preFS embed.FS

//go:embed data/post/*.txt
var postFS embed.FS

//go:embed data/strip/*.csv
var stripFS embed.FS

//go:embed data/puncnorm.csv
var puncnormData []byte

package epitran

import "errors"

// MappingError indicates a problem with grapheme-to-phoneme mapping data.
type MappingError struct {
	msg string
}

func (e *MappingError) Error() string { return e.msg }

// DatafileError indicates a problem with a data file.
type DatafileError struct {
	msg string
}

func (e *DatafileError) Error() string { return e.msg }

// ErrUnsupportedSpecialLanguage is returned for language codes that require
// special handlers not yet ported (eng-Latn, cmn-Hans, cmn-Hant, jpn-Jpan, yue-Hant).
var ErrUnsupportedSpecialLanguage = errors.New("epitran: this language code requires a special handler that is not yet implemented")

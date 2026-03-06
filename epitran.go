package epitran

// specialCodes are language codes that require special handlers not yet ported.
var specialCodes = map[string]bool{
	"eng-Latn": true,
	"cmn-Hans": true,
	"cmn-Hant": true,
	"jpn-Jpan": true,
	"yue-Hant": true,
}

// Option configures an Epitran instance.
type Option func(*config)

type config struct {
	preproc   bool
	postproc  bool
	ligatures bool
	tones     bool
}

func defaultConfig() config {
	return config{
		preproc:  true,
		postproc: true,
	}
}

// WithPreproc enables or disables preprocessing rules.
func WithPreproc(v bool) Option { return func(c *config) { c.preproc = v } }

// WithPostproc enables or disables postprocessing rules.
func WithPostproc(v bool) Option { return func(c *config) { c.postproc = v } }

// WithLigatures enables or disables affricate ligatures in output.
func WithLigatures(v bool) Option { return func(c *config) { c.ligatures = v } }

// WithTones enables or disables tone marks in output.
func WithTones(v bool) Option { return func(c *config) { c.tones = v } }

// Epitran provides grapheme-to-phoneme transliteration.
type Epitran struct {
	simple *simpleEpitran
}

// New creates an Epitran instance for the given language-script code.
// Returns ErrUnsupportedSpecialLanguage for codes requiring special handlers
// (eng-Latn, cmn-Hans, cmn-Hant, jpn-Jpan, yue-Hant).
func New(code string, opts ...Option) (*Epitran, error) {
	if specialCodes[code] {
		return nil, ErrUnsupportedSpecialLanguage
	}

	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	se, err := newSimpleEpitran(code, cfg.preproc, cfg.postproc, cfg.ligatures, cfg.tones)
	if err != nil {
		return nil, err
	}

	return &Epitran{simple: se}, nil
}

// Transliterate converts orthographic text to IPA.
// Unmapped characters are passed through to the output.
func (e *Epitran) Transliterate(text string) string {
	return e.simple.transliterate(text, false)
}

// TransliterateNorm converts orthographic text to IPA with punctuation normalization.
func (e *Epitran) TransliterateNorm(text string) string {
	return e.simple.transliterate(text, true)
}

// StrictTrans converts orthographic text to IPA, dropping unmapped characters.
func (e *Epitran) StrictTrans(text string) string {
	return e.simple.strictTrans(text, false)
}

// StrictTransNorm converts orthographic text to IPA with punctuation normalization,
// dropping unmapped characters.
func (e *Epitran) StrictTransNorm(text string) string {
	return e.simple.strictTrans(text, true)
}

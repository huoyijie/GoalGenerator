package model

type Lazy interface {
	Lazy()
}

type Purge interface {
	Purge()
}

type Ctrl interface {
	Ctrl()
}

type Translate interface {
	TranslatePkg() map[string]string

	TranslateName() map[string]string

	TranslateFields() map[string]map[string]string

	TranslateOptions() map[string]map[string]map[string]string
}

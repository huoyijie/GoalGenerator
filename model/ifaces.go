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

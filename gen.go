package goalgenerator

import (
	"gorm.io/gorm"
)

type ILazy interface {
	IsLazy() bool
}

type Lazy struct{}

func (*Lazy) IsLasy() bool {
	return true
}

type Base gorm.Model

type IValid interface {
	Valid() error
}

type IProp interface {
	IsProp() bool
	Prop() (string, string)
}

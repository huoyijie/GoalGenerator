package goalgenerator

import (
	"gorm.io/gorm"
	"time"
)

type ILazy interface {
	Lazy()
}

type Lazy struct{}

// Lazy implements ILazy
func (*Lazy) Lazy() {}

var _ ILazy = (*Lazy)(nil)

type Base struct {
	ID        uint      `gorm:"primarykey" goal:"<number>primary,sortable,asc,uint"`
	CreatedAt time.Time `goal:"<calendar>autowired"`
	UpdatedAt time.Time `goal:"<calendar>autowired"`
	// todo provide <autowired> component
	DeletedAt gorm.DeletedAt `gorm:"index" goal:"<calendar>autowired"`
}

type IValid interface {
	Valid() error
}

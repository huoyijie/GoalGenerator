package goalgenerator

import (
	"gorm.io/gorm"
	"time"
)

type Lazy interface {
	Lazy()
}

type Base struct {
	ID        uint           `gorm:"primarykey" goal:"<number>primary,sortable,asc,uint"`
	CreatedAt time.Time      `goal:"<calendar>autowired"`
	UpdatedAt time.Time      `goal:"<calendar>autowired"`
	DeletedAt gorm.DeletedAt `gorm:"index" goal:"<calendar>autowired"`
	Creator   uint           `gorm:"index" goal:"<number>autowired,uint"`
}

type IValid interface {
	Valid() error
}

package goalgenerator

import (
	"errors"

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

type Field struct {
	Name           string          `yaml:",omitempty"`
	StorageRules   []StorageRule   `yaml:"storageRules,omitempty"`
	Component      Component       `yaml:",omitempty"`
	ComponentRules []ComponentRule `yaml:"componentRules,omitempty"`
	ValidateRules  []ValidateRule  `yaml:"validateRules,omitempty"`
}

// Valid implements IValid
func (f *Field) Valid() error {
	if f.Name == "" {
		return errors.New("field's name is required")
	}

	for _, r := range f.StorageRules {
		if err := r.ValidField(); err != nil {
			return err
		}
	}

	if err := f.Component.Valid(); err != nil {
		return err
	}

	for _, r := range f.ComponentRules {
		if err := r.ValidField(); err != nil {
			return err
		}
	}

	for _, r := range f.ValidateRules {
		if err := r.Valid(); err != nil {
			return err
		}
	}
	return nil
}

var _ IValid = (*Field)(nil)

type Model struct {
	Name           string          `yaml:",omitempty"`
	StorageRules   []StorageRule   `yaml:"storageRules,omitempty"`
	ComponentRules []ComponentRule `yaml:"componentRules,omitempty"`
	Fields         []Field         `yaml:",omitempty"`
}

// Valid implements IValid
func (m *Model) Valid() error {
	if m.Name == "" {
		return errors.New("model's name is required")
	}

	for _, r := range m.StorageRules {
		if err := r.ValidModel(); err != nil {
			return err
		}
	}

	for _, r := range m.ComponentRules {
		if err := r.ValidModel(); err != nil {
			return err
		}
	}

	for _, f := range m.Fields {
		if err := f.Valid(); err != nil {
			return err
		}
	}

	return nil
}

var _ IValid = (*Model)(nil)

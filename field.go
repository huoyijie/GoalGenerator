package goalgenerator

import "errors"

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

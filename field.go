package goalgenerator

import (
	"errors"
	"strings"
)

type Field struct {
	Name           string          `yaml:",omitempty"`
	StorageRules   []StorageRule   `yaml:"storageRules,omitempty"`
	Component      Component       `yaml:",omitempty"`
	ComponentRules []ComponentRule `yaml:"componentRules,omitempty"`
	ValidateRules  []ValidateRule  `yaml:"validateRules,omitempty"`
}

func (f *Field) Type() (t string) {
	switch f.Component {
	case COMPONENT_NUMBER:
		t = "int"
		for _, r := range f.ComponentRules {
			if r == COMPONENT_RULE_FLOAT {
				t = "float"
				break
			}
		}
	case COMPONENT_UUID, COMPONENT_TEXT, COMPONENT_PASSWORD, COMPONENT_FILE:
		t = "string"
	case COMPONENT_CALENDAR:
		t = "time.Time"
	case COMPONENT_SWITCH:
		t = "bool"
	case COMPONENT_DROPDOWN:
		for _, r := range f.ComponentRules {
			if r.IsProp() {
				if k, v := r.Prop(); ComponentRule(k) == COMPONENT_RULE_BELONGTO {
					parts := strings.Split(v, ".")
					t = strings.Join(parts[:len(parts)-1], ".")
					break
				}
			}
		}
	}
	return
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

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
			} else if r == COMPONENT_RULE_UINT {
				t = "uint"
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

func (f *Field) Tag() (tag string) {
	sb := strings.Builder{}
	var primary, unique bool
	if len(f.StorageRules) > 0 {
		sb.WriteString(`gorm:"`)
		for i, r := range f.StorageRules {
			if r == STORAGE_RULE_PRIMARY {
				primary = true
			} else if r == STORAGE_RULE_UNIQUE {
				unique = true
			}
			if r.IsProp() {
				sb.WriteString(string(r)[1:])
			} else {
				sb.WriteString(string(r))
			}
			if i < len(f.StorageRules)-1 {
				sb.WriteRune(',')
			} else {
				sb.WriteString(`" `)
			}
		}
	}

	if len(f.ValidateRules) > 0 {
		sb.WriteString(`binding:"`)
		for i, r := range f.ValidateRules {
			if r.IsProp() {
				sb.WriteString(string(r)[1:])
			} else {
				sb.WriteString(string(r))
			}
			if i < len(f.ValidateRules)-1 {
				sb.WriteRune(',')
			} else {
				sb.WriteString(`" `)
			}
		}
	}

	sb.WriteString(`goal:"<`)
	sb.WriteString(string(f.Component))
	sb.WriteRune('>')
	if primary {
		sb.WriteString("primary,")
	}
	if unique {
		sb.WriteString("unique,")
	}
	for i, r := range f.ComponentRules {
		if r.IsProp() {
			sb.WriteString(string(r)[1:])
		} else {
			sb.WriteString(string(r))
		}
		if i < len(f.ComponentRules)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteString(`"`)

	return sb.String()
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

	isNumber := f.Component == COMPONENT_NUMBER
	var isUint, isFloat bool
	for _, r := range f.ComponentRules {
		if err := r.ValidField(); err != nil {
			return err
		}
		if isNumber {
			if r == COMPONENT_RULE_FLOAT {
				isFloat = true
			} else if r == COMPONENT_RULE_UINT {
				isUint = true
			}
		}
	}
	if isUint && isFloat {
		return errors.New("number can not be set uint and float at the same time")
	}

	for _, r := range f.ValidateRules {
		if err := r.Valid(); err != nil {
			return err
		}
	}
	return nil
}

var _ IValid = (*Field)(nil)

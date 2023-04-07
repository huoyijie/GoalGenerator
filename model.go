package goalgenerator

import (
	"errors"
	"fmt"
	"strings"
)

type Model struct {
	Package,
	Name string `yaml:",omitempty"`
	StorageRules   []StorageRule   `yaml:"storageRules,omitempty"`
	ComponentRules []ComponentRule `yaml:"componentRules,omitempty"`
	Fields         []Field         `yaml:",omitempty"`
}

func (m *Model) Version() string {
	return Version
}

func (m *Model) Imports() (imports []string) {
	if m.EmbeddingBase() || m.Lazy() {
		imports = append(imports, fmt.Sprintf(`"%s"`, GetMoudlePath()))
	}
	for _, f := range m.Fields {
		if t := f.Type(); t == "time.Time" {
			imports = append(imports, `"time"`)
		} else if parts := strings.Split(t, "."); len(parts) == 2 {
			imports = append(imports, fmt.Sprintf(`"%s/%s"`, GetMoudlePath(), parts[0]))
		}
	}
	return
}

func (m *Model) EmbeddingBase() bool {
	for _, r := range m.StorageRules {
		if r == STORAGE_RULE_EMBEDDING_BASE {
			return true
		}
	}
	return false
}

func (m *Model) Lazy() bool {
	for _, r := range m.ComponentRules {
		if r == COMPONENT_RULE_LAZY {
			return true
		}
	}
	return false
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

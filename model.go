package goalgenerator

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const Version string = "0.0.12"

//go:embed template/*.tpl
var tmpl string

//go:embed goal-schema.json
var schema string

type ILazy interface {
	Lazy()
}

type Lazy struct{}

// Lazy implements ILazy
func (*Lazy) Lazy() {}

var _ ILazy = (*Lazy)(nil)

type Base = gorm.Model

type IValid interface {
	Valid() error
}

func GetMoudlePath() (pkgPath string) {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}

	pkgPath = modfile.ModulePath(goModBytes)
	return
}

type Model struct {
	Package,
	Name string `yaml:",omitempty"`
	Database *struct {
		EmbeddingBase bool `yaml:",omitempty"`
	} `yaml:",omitempty"`
	View *struct {
		Lazy bool `yaml:",omitempty"`
	} `yaml:",omitempty"`
	Fields []Field `yaml:",omitempty"`
}

func (m *Model) Version() string {
	return Version
}

func (m *Model) Gen() error {
	for i := range m.Fields {
		m.Fields[i].Model = m
	}

	os.Mkdir(m.Package, os.ModePerm)

	f, err := os.Create(fmt.Sprintf("%s/%s.go", m.Package, m.Name))
	if err != nil {
		return err
	}

	t := template.Must(template.New("").Parse(tmpl))

	return t.Execute(f, m)
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
	return m.Database != nil && m.Database.EmbeddingBase
}

func (m *Model) Lazy() bool {
	return m.View != nil && m.View.Lazy
}

// Valid implements IValid
func (m *Model) Valid() error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	var v any
	if err := yaml.Unmarshal(data, &v); err != nil {
		return err
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("goal-schema.json", strings.NewReader(schema)); err != nil {
		return err
	}

	sch, err := c.Compile("goal-schema.json")
	if err != nil {
		return err
	}

	if err = sch.Validate(v); err != nil {
		return err
	}

	return nil
}

var _ IValid = (*Model)(nil)

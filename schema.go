package goalgenerator

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

const Version string = "0.0.21"

//go:embed template/*.tpl
var tmpl string

//go:embed schema.json
var schema string

//go:embed go.mod
var goMod []byte

type IValid interface {
	Valid() error
}

func ToLowerFirstLetter(str string) string {
	a := []rune(str)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func GetMoudlePath() (pkgPath string) {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}

	pkgPath = modfile.ModulePath(goModBytes)
	return
}

type Translate struct {
	En    string `yaml:",omitempty"`
	Zh_CN string `yaml:"zh_CN,omitempty"`
}

type Model struct {
	Package *struct {
		Value     string `yaml:",omitempty"`
		Translate `yaml:",inline,omitempty"`
	} `yaml:",omitempty"`
	Name *struct {
		Value     string `yaml:",omitempty"`
		Translate `yaml:",inline,omitempty"`
	} `yaml:",omitempty"`
	Database *struct {
		EmbeddingBase,
		Purge bool `yaml:",omitempty"`
		TableName string `yaml:",omitempty"`
	} `yaml:",omitempty"`
	View *struct {
		Lazy,
		Ctrl bool `yaml:",omitempty"`
		Icon string `yaml:",omitempty"`
	} `yaml:",omitempty"`
	Fields []Field `yaml:",omitempty"`
}

func (m *Model) Dropdowns() (fields []Field) {
	for _, f := range m.Fields {
		if f.DropdownStrings() || f.DropdownInts() || f.DropdownUints() || f.DropdownFloats() || f.DropdownDynamicStrings() || f.DropdownDynamicInts() || f.DropdownDynamicUints() || f.DropdownDynamicFloats() {
			fields = append(fields, f)
		}
	}
	return
}

func (m *Model) Version() string {
	return Version
}

func (m *Model) Gen() error {
	for i := range m.Fields {
		m.Fields[i].Model = m
	}

	os.Mkdir(m.Package.Value, os.ModePerm)

	f, err := os.Create(fmt.Sprintf("%s/%s.go", m.Package.Value, m.Name.Value))
	if err != nil {
		return err
	}

	t := template.Must(template.New("").Parse(tmpl))

	return t.Execute(f, m)
}

func (m *Model) Imports() (imports []string) {
	if m.EmbeddingBase() {
		imports = append(imports, fmt.Sprintf(`"%s/model"`, modfile.ModulePath(goMod)))
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

func (m *Model) CustomTableName() bool {
	return m.Database != nil && m.Database.TableName != ""
}

func (m *Model) TableName() string {
	return m.Database.TableName
}

func (m *Model) Purge() bool {
	return m.Database != nil && m.Database.Purge
}

func (m *Model) Lazy() bool {
	return m.View.Lazy
}

func (m *Model) Ctrl() bool {
	return m.View.Ctrl
}

func (m *Model) Icon() string {
	return m.View.Icon
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
	if err := c.AddResource("schema.json", strings.NewReader(schema)); err != nil {
		return err
	}

	sch, err := c.Compile("schema.json")
	if err != nil {
		return err
	}

	if err = sch.Validate(v); err != nil {
		return err
	}

	return nil
}

var _ IValid = (*Model)(nil)

package goalgenerator

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"golang.org/x/mod/modfile"
	"gorm.io/gorm"
)

const Version string = "0.0.5"

type ILazy interface {
	IsLazy() bool
}

type Lazy struct{}

func (*Lazy) IsLasy() bool {
	return true
}

type Base = gorm.Model

type IValid interface {
	Valid() error
}

type IProp interface {
	IsProp() bool
	Prop() (string, string)
}

//go:embed template/*.tpl
var tmpl string

func GenModel(m *Model) error {
	os.Mkdir(m.Package, os.ModePerm)

	f, err := os.Create(fmt.Sprintf("%s/%s.go", m.Package, m.Name))
	if err != nil {
		return err
	}

	t := template.Must(template.New("").Parse(tmpl))

	return t.Execute(f, m)
}

func GetMoudlePath() (pkgPath string) {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}

	pkgPath = modfile.ModulePath(goModBytes)
	return
}

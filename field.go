package goalgenerator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func ToLowerFirstLetter(str string) string {
	a := []rune(str)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

type Field struct {
	Model    *Model `yaml:"-"`
	Name     string `yaml:",omitempty"`
	Database *struct {
		PrimaryKey,
		Unique,
		Index bool `yaml:",omitempty"`
	} `yaml:",omitempty"`
	View *struct {
		Name,
		BelongTo,
		UploadTo string `yaml:",omitempty"`
		Sortable,
		Asc,
		Desc,
		GlobalSearch,
		Secret,
		Hidden,
		Readonly,
		Postonly,
		Autowired,
		Filter,
		ShowTime,
		ShowIcon,
		Uint,
		Float bool `yaml:",omitempty"`
	} `yaml:",omitempty"`
	Validator *struct {
		Required,
		Email,
		Alphanum,
		Alpha bool `yaml:",omitempty"`
		Min,
		Max,
		Len *int `yaml:",omitempty"`
	} `yaml:",omitempty"`
}

func (f *Field) Type() (t string) {
	switch f.View.Name {
	case "number":
		if f.View.Float {
			t = "float"
		} else if f.View.Uint {
			t = "uint"
		} else {
			t = "int"
		}
	case "uuid", "text", "password", "file":
		t = "string"
	case "calendar":
		t = "time.Time"
	case "switch":
		t = "bool"
	case "dropdown":
		parts := strings.Split(f.View.BelongTo, ".")
		t = strings.Join(parts[:len(parts)-1], ".")
	}
	return
}

func (f *Field) Tag() (tag string) {
	sb := strings.Builder{}
	var primary, unique bool
	if f.Database != nil {
		primary = f.Database.PrimaryKey
		unique = f.Database.Unique
		sb.WriteString(`gorm:"`)
		var hasPrev bool
		if primary {
			if hasPrev {
				sb.WriteRune(',')
			}
			sb.WriteString("primarykey")
			hasPrev = true
		}
		if unique {
			if hasPrev {
				sb.WriteRune(',')
			}
			sb.WriteString("unique")
			hasPrev = true
		}
		if f.Database.Index {
			if hasPrev {
				sb.WriteRune(',')
			}
			sb.WriteString("index")
		}
		sb.WriteString(`" `)
	}

	if f.Validator != nil {
		sb.WriteString(`binding:"`)
		t := reflect.TypeOf(f.Validator).Elem()
		val := reflect.ValueOf(f.Validator).Elem()
		var hasPrev bool
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fVal := val.FieldByName(f.Name)
			switch fVal.Kind() {
			case reflect.Bool:
				if fVal.Bool() {
					if hasPrev {
						sb.WriteRune(',')
					}
					sb.WriteString(ToLowerFirstLetter(f.Name))
					hasPrev = true
				}
			case reflect.Pointer:
				if !fVal.IsNil() {
					fVal = fVal.Elem()
					if hasPrev {
						sb.WriteRune(',')
					}
					sb.WriteString(ToLowerFirstLetter(f.Name))
					sb.WriteRune('=')
					sb.WriteString(fmt.Sprintf("%v", fVal.Interface()))
					hasPrev = true
				}
			}
		}
		sb.WriteString(`" `)
	}

	if f.View != nil {
		sb.WriteString(`goal:"<`)
		sb.WriteString(f.View.Name)
		sb.WriteRune('>')

		var hasPrev bool
		if primary {
			if hasPrev {
				sb.WriteRune(',')
			}
			sb.WriteString("primary")
			hasPrev = true
		}
		if unique {
			if hasPrev {
				sb.WriteRune(',')
			}
			sb.WriteString("unique")
			hasPrev = true
		}

		pkg := f.Model.Package
		t := reflect.TypeOf(f.View).Elem()
		val := reflect.ValueOf(f.View).Elem()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Name == "Name" {
				continue
			}
			fVal := val.FieldByName(f.Name)
			switch fVal.Kind() {
			case reflect.Bool:
				if fVal.Bool() {
					if hasPrev {
						sb.WriteRune(',')
					}
					sb.WriteString(ToLowerFirstLetter(f.Name))
					hasPrev = true
				}
			case reflect.String:
				if !fVal.IsZero() {
					if hasPrev {
						sb.WriteRune(',')
					}
					sb.WriteString(ToLowerFirstLetter(f.Name))
					sb.WriteRune('=')
					if f.Name == "BelongTo" && strings.Count(fVal.String(), ".") == 1 {
						sb.WriteString(pkg)
						sb.WriteRune('.')
					}
					sb.WriteString(fVal.String())
					hasPrev = true
				}
			}
		}
		sb.WriteString(`"`)
	}

	return sb.String()
}

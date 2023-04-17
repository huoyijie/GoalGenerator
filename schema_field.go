package goalgenerator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Field struct {
	Model *Model `yaml:"-"`
	Name  *struct {
		Value     string `yaml:",omitempty"`
		Translate `yaml:",inline,omitempty"`
	} `yaml:",omitempty"`
	Database *struct {
		PrimaryKey,
		Unique,
		Index bool `yaml:",omitempty"`
	} `yaml:",omitempty"`
	View *struct {
		Base *struct {
			Readonly,
			Postonly,
			Sortable,
			Asc,
			Desc,
			GlobalSearch,
			Filter,
			Hidden,
			Secret,
			Autowired bool `yaml:",omitempty"`
		} `yaml:",omitempty"`
		Uuid,
		Text,
		Switch,
		Password bool `yaml:",omitempty"`
		Number *struct {
			ShowButtons,
			Float,
			Uint bool `yaml:",omitempty"`
			Min,
			Max *int `yaml:",omitempty"`
		} `yaml:",omitempty"`
		File *struct {
			UploadTo string `yaml:",omitempty"`
		} `yaml:",omitempty"`
		Dropdown *struct {
			DynamicStrings,
			DynamicInts,
			DynamicUints,
			DynamicFloats bool `yaml:",omitempty"`
			Strings  []string  `yaml:",omitempty"`
			Ints     []int     `yaml:",omitempty"`
			Uints    []uint    `yaml:",omitempty"`
			Floats   []float64 `yaml:",omitempty"`
			BelongTo *struct {
				Pkg,
				Name,
				Field string `yaml:",omitempty"`
			} `yaml:",omitempty"`
		} `yaml:",omitempty"`
		Calendar *struct {
			ShowTime,
			ShowIcon bool `yaml:",omitempty"`
		} `yaml:",omitempty"`
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

func (f *Field) Component() (c string) {
	switch {
	case f.View.Number != nil:
		c = "number"
	case f.View.Uuid:
		c = "uuid"
	case f.View.Text:
		c = "text"
	case f.View.Password:
		c = "password"
	case f.View.File != nil:
		c = "file"
	case f.View.Calendar != nil:
		c = "calendar"
	case f.View.Switch:
		c = "switch"
	case f.View.Dropdown != nil:
		c = "dropdown"
	}
	return
}

func (f *Field) Type() (t string) {
	switch {
	case f.View.Number != nil:
		if f.View.Number.Float {
			t = "float64"
		} else if f.View.Number.Uint {
			t = "uint"
		} else {
			t = "int"
		}
	case f.View.Uuid, f.View.Text, f.View.Password, f.View.File != nil:
		t = "string"
	case f.View.Calendar != nil:
		t = "time.Time"
	case f.View.Switch:
		t = "bool"
	case f.View.Dropdown != nil:
		switch {
		case f.View.Dropdown.BelongTo != nil:
			if belongTo := f.View.Dropdown.BelongTo; belongTo.Pkg == "" {
				t = belongTo.Name
			} else {
				t = strings.Join([]string{belongTo.Pkg, belongTo.Name}, ".")
			}
		case f.View.Dropdown.DynamicStrings, f.View.Dropdown.Strings != nil:
			t = "string"
		case f.View.Dropdown.DynamicInts, f.View.Dropdown.Ints != nil:
			t = "int"
		case f.View.Dropdown.DynamicUints, f.View.Dropdown.Uints != nil:
			t = "uint"
		case f.View.Dropdown.DynamicFloats, f.View.Dropdown.Floats != nil:
			t = "float64"
		}
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
		sb.WriteString(f.Component())
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

		m := f.Model
		t := reflect.TypeOf(f.View).Elem()
		val := reflect.ValueOf(f.View).Elem()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fVal := val.FieldByName(f.Name)
			switch fVal.Kind() {
			case reflect.Pointer:
				if !fVal.IsNil() {
					ct := fVal.Type().Elem()
					cVal := fVal.Elem()
					for i := 0; i < ct.NumField(); i++ {
						cf := ct.Field(i)
						cfVal := cVal.FieldByName(cf.Name)
						switch cfVal.Kind() {
						case reflect.Bool:
							if cfVal.Bool() {
								if hasPrev {
									sb.WriteRune(',')
								}
								sb.WriteString(ToLowerFirstLetter(cf.Name))
								hasPrev = true
							}
						case reflect.String:
							if !cfVal.IsZero() {
								if hasPrev {
									sb.WriteRune(',')
								}
								sb.WriteString(ToLowerFirstLetter(cf.Name))
								sb.WriteRune('=')
								sb.WriteString(cfVal.String())
								hasPrev = true
							}
						case reflect.Pointer:
							if !cfVal.IsNil() {
								switch e := cfVal.Elem(); e.Kind() {
								case reflect.Int:
									if hasPrev {
										sb.WriteRune(',')
									}
									sb.WriteString(ToLowerFirstLetter(cf.Name))
									sb.WriteRune('=')
									sb.WriteString(fmt.Sprintf("%v", e.Interface()))
									hasPrev = true
								case reflect.Struct:
									if f.Name == "Dropdown" && cf.Name == "BelongTo" {
										p := e.FieldByName("Pkg").Interface().(string)
										if p == "" {
											p = m.Package.Value
										}
										n := e.FieldByName("Name").Interface()
										fn := e.FieldByName("Field").Interface()
										if hasPrev {
											sb.WriteRune(',')
										}
										sb.WriteString(ToLowerFirstLetter(cf.Name))
										sb.WriteRune('=')
										sb.WriteString(fmt.Sprintf("%s.%s.%s", p, n, fn))
										hasPrev = true
									}
								}
							}
						case reflect.Slice:
							if !cfVal.IsZero() && cfVal.Len() > 0 {
								if f.Name == "Dropdown" {
									if hasPrev {
										sb.WriteRune(',')
									}
									sb.WriteString(ToLowerFirstLetter(cf.Name))
									//cfVal.Interface()
									hasPrev = true
								}
							}
						}
					}
				}
			}
		}
		sb.WriteString(`"`)
	}

	return sb.String()
}

func (f *Field) DropdownStrings() bool {
	return f.View.Dropdown != nil && len(f.View.Dropdown.Strings) > 0
}

func (f *Field) OptionStrings() string {
	if f.DropdownStrings() {
		sb := strings.Builder{}
		sb.WriteString("[]string{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Strings {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(option)
			sb.WriteRune('"')
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownInts() bool {
	return f.View.Dropdown != nil && len(f.View.Dropdown.Ints) > 0
}

func (f *Field) OptionInts() string {
	if f.DropdownInts() {
		sb := strings.Builder{}
		sb.WriteString("[]int{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Ints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(fmt.Sprintf("%d", option))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownUints() bool {
	return f.View.Dropdown != nil && len(f.View.Dropdown.Uints) > 0
}

func (f *Field) OptionUints() string {
	if f.DropdownUints() {
		sb := strings.Builder{}
		sb.WriteString("[]uint{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Uints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(fmt.Sprintf("%d", option))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownFloats() bool {
	return f.View.Dropdown != nil && len(f.View.Dropdown.Floats) > 0
}

func (f *Field) OptionFloats() string {
	if f.DropdownFloats() {
		sb := strings.Builder{}
		sb.WriteString("[]float64{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Floats {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(strconv.FormatFloat(option, 'f', -1, 64))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownDynamicStrings() bool {
	return f.View.Dropdown != nil && f.View.Dropdown.DynamicStrings
}

func (f *Field) DropdownDynamicInts() bool {
	return f.View.Dropdown != nil && f.View.Dropdown.DynamicInts
}

func (f *Field) DropdownDynamicUints() bool {
	return f.View.Dropdown != nil && f.View.Dropdown.DynamicUints
}

func (f *Field) DropdownDynamicFloats() bool {
	return f.View.Dropdown != nil && f.View.Dropdown.DynamicFloats
}

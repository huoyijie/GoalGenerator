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
			Option *struct {
				Dynamic *struct {
					Strings,
					Ints,
					Uints,
					Floats bool `yaml:",omitempty"`
				} `yaml:",omitempty"`
				Strings []struct {
					Value     string `yaml:",omitempty"`
					Translate `yaml:",inline,omitempty"`
				} `yaml:",omitempty"`
				Ints []struct {
					Value     *int `yaml:",omitempty"`
					Translate `yaml:",inline,omitempty"`
				} `yaml:",omitempty"`
				Uints []struct {
					Value     *uint `yaml:",omitempty"`
					Translate `yaml:",inline,omitempty"`
				} `yaml:",omitempty"`
				Floats []struct {
					Value     *float64 `yaml:",omitempty"`
					Translate `yaml:",inline,omitempty"`
				} `yaml:",omitempty"`
			} `yaml:",omitempty"`
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
	t := reflect.TypeOf(f.View).Elem()
	v := reflect.ValueOf(f.View).Elem()
	for i := 0; i < t.NumField(); i++ {
		vf := t.Field(i)
		if vf.Name == "Base" {
			continue
		}
		vfVal := v.FieldByName(vf.Name)
		if !vfVal.IsZero() {
			c = ToLowerFirstLetter(vf.Name)
		}
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
		case f.DropdownStrings(), f.DropdownDynamicStrings():
			t = "string"
		case f.DropdownInts(), f.DropdownDynamicInts():
			t = "int"
		case f.DropdownUints(), f.DropdownDynamicUints():
			t = "uint"
		case f.DropdownFloats(), f.DropdownDynamicFloats():
			t = "float64"
		}
	}
	return
}

func (f *Field) writeString(sb *strings.Builder, hasPrev *bool, value string, more ...string) {
	if *hasPrev {
		sb.WriteRune(',')
	} else {
		*hasPrev = true
	}
	sb.WriteString(value)
	for _, s := range more {
		sb.WriteString(s)
	}
}

func (f *Field) gorm(sb *strings.Builder) (primary bool, unique bool) {
	if f.Database != nil {
		primary = f.Database.PrimaryKey
		unique = f.Database.Unique
		sb.WriteString(`gorm:"`)
		var hasPrev bool
		if primary {
			f.writeString(sb, &hasPrev, "primarykey")
		}
		if unique {
			f.writeString(sb, &hasPrev, "unique")
		}
		if f.Database.Index {
			f.writeString(sb, &hasPrev, "index")
		}
		sb.WriteString(`" `)
	}
	return
}

func (f *Field) validator(sb *strings.Builder) {
	if f.Validator != nil {
		sb.WriteString(`binding:"`)
		t := reflect.TypeOf(f.Validator).Elem()
		val := reflect.ValueOf(f.Validator).Elem()
		var hasPrev bool
		for i := 0; i < t.NumField(); i++ {
			vf := t.Field(i)
			fVal := val.FieldByName(vf.Name)
			switch fVal.Kind() {
			case reflect.Bool:
				if fVal.Bool() {
					f.writeString(sb, &hasPrev, ToLowerFirstLetter(vf.Name))
				}
			case reflect.Pointer:
				if !fVal.IsNil() {
					fVal = fVal.Elem()
					f.writeString(sb, &hasPrev, ToLowerFirstLetter(vf.Name), "=", fmt.Sprintf("%v", fVal.Interface()))
				}
			}
		}
		sb.WriteString(`" `)
	}
}

func (f *Field) view(sb *strings.Builder, primary, unique bool) {
	if f.View != nil {
		sb.WriteString(`goal:"<`)
		sb.WriteString(f.Component())
		sb.WriteRune('>')

		var hasPrev bool
		if primary {
			f.writeString(sb, &hasPrev, "primary")
		}
		if unique {
			f.writeString(sb, &hasPrev, "unique")
		}

		m := f.Model
		t := reflect.TypeOf(f.View).Elem()
		val := reflect.ValueOf(f.View).Elem()
		for i := 0; i < t.NumField(); i++ {
			vf := t.Field(i)
			fVal := val.FieldByName(vf.Name)
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
								f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name))
							}
						case reflect.String:
							if !cfVal.IsZero() {
								f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name), "=", cfVal.String())
							}
						case reflect.Pointer:
							if !cfVal.IsNil() {
								switch e := cfVal.Elem(); e.Kind() {
								case reflect.Int:
									f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name), "=", fmt.Sprintf("%v", e.Interface()))
								case reflect.Struct:
									if vf.Name == "Dropdown" {
										switch cf.Name {
										case "BelongTo":
											p := e.FieldByName("Pkg").Interface().(string)
											if p == "" {
												p = m.Package.Value
											}
											n := e.FieldByName("Name").Interface()
											fn := e.FieldByName("Field").Interface()
											f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name), "=", fmt.Sprintf("%s.%s.%s", p, n, fn))
										case "Option":
											ot := reflect.TypeOf(e.Interface())
											for i := 0; i < ot.NumField(); i++ {
												of := ot.Field(i)
												oVal := e.FieldByName(of.Name)
												switch oVal.Kind() {
												case reflect.Slice:
													if !oVal.IsZero() && oVal.Len() > 0 {
														f.writeString(sb, &hasPrev, ToLowerFirstLetter(of.Name))
													}
												case reflect.Pointer: // dynamic
													if !oVal.IsNil() {
														dt := reflect.TypeOf(oVal.Elem().Interface())
														for i := 0; i < dt.NumField(); i++ {
															df := dt.Field(i)
															dVal := oVal.Elem().FieldByName(df.Name)
															if dVal.Bool() {
																f.writeString(sb, &hasPrev, "dynamic", df.Name)
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		sb.WriteString(`"`)
	}
}

func (f *Field) Tag() (tag string) {
	sb := strings.Builder{}
	primary, unique := f.gorm(&sb)
	f.validator(&sb)
	f.view(&sb, primary, unique)
	return sb.String()
}

func (f *Field) DropdownStrings() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && len(d.Option.Strings) > 0 {
		return true
	}
	return false
}

func (f *Field) OptionStrings() string {
	if f.DropdownStrings() {
		sb := strings.Builder{}
		sb.WriteString("[]string{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Strings {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(option.Value)
			sb.WriteRune('"')
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) OptionStringLabels() string {
	if f.DropdownStrings() {
		sb := strings.Builder{}
		sb.WriteString("map[string]map[string]string{")
		sb.WriteString(` "en": {`)
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Strings {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(option.Value)
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.En)
			sb.WriteRune('"')
		}
		sb.WriteString(`}, "zh_CN": {`)

		hasPrev = false
		for _, option := range f.View.Dropdown.Option.Strings {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(option.Value)
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.Zh_CN)
			sb.WriteRune('"')
		}
		sb.WriteString("}")

		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownInts() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && len(d.Option.Ints) > 0 {
		return true
	}
	return false
}

func (f *Field) OptionInts() string {
	if f.DropdownInts() {
		sb := strings.Builder{}
		sb.WriteString("[]int{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Ints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) OptionIntLabels() string {
	if f.DropdownInts() {
		sb := strings.Builder{}
		sb.WriteString("map[string]map[string]string{")
		sb.WriteString(` "en": {`)
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Ints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.En)
			sb.WriteRune('"')
		}
		sb.WriteString(`}, "zh_CN": {`)

		hasPrev = false
		for _, option := range f.View.Dropdown.Option.Ints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.Zh_CN)
			sb.WriteRune('"')
		}
		sb.WriteString("}")

		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownUints() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && len(d.Option.Uints) > 0 {
		return true
	}
	return false
}

func (f *Field) OptionUints() string {
	if f.DropdownUints() {
		sb := strings.Builder{}
		sb.WriteString("[]uint{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Uints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) OptionUintLabels() string {
	if f.DropdownUints() {
		sb := strings.Builder{}
		sb.WriteString("map[string]map[string]string{")
		sb.WriteString(` "en": {`)
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Uints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.En)
			sb.WriteRune('"')
		}
		sb.WriteString(`}, "zh_CN": {`)

		hasPrev = false
		for _, option := range f.View.Dropdown.Option.Uints {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(fmt.Sprintf("%d", *option.Value))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.Zh_CN)
			sb.WriteRune('"')
		}
		sb.WriteString("}")

		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownFloats() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && len(d.Option.Floats) > 0 {
		return true
	}
	return false
}

func (f *Field) OptionFloats() string {
	if f.DropdownFloats() {
		sb := strings.Builder{}
		sb.WriteString("[]float64{")
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Floats {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteString(strconv.FormatFloat(*option.Value, 'f', -1, 64))
		}
		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) OptionFloatLabels() string {
	if f.DropdownFloats() {
		sb := strings.Builder{}
		sb.WriteString("map[string]map[string]string{")
		sb.WriteString(` "en": {`)
		var hasPrev bool
		for _, option := range f.View.Dropdown.Option.Floats {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(strconv.FormatFloat(*option.Value, 'f', -1, 64))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.En)
			sb.WriteRune('"')
		}
		sb.WriteString(`}, "zh_CN": {`)

		hasPrev = false
		for _, option := range f.View.Dropdown.Option.Floats {
			if hasPrev {
				sb.WriteString(", ")
			} else {
				hasPrev = true
			}
			sb.WriteRune('"')
			sb.WriteString(strconv.FormatFloat(*option.Value, 'f', -1, 64))
			sb.WriteRune('"')
			sb.WriteString(": ")
			sb.WriteRune('"')
			sb.WriteString(option.Zh_CN)
			sb.WriteRune('"')
		}
		sb.WriteString("}")

		sb.WriteString("}")
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownDynamicStrings() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && d.Option.Dynamic != nil && d.Option.Dynamic.Strings {
		return true
	}
	return false
}

func (f *Field) DropdownDynamicInts() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && d.Option.Dynamic != nil && d.Option.Dynamic.Ints {
		return true
	}
	return false
}

func (f *Field) DropdownDynamicUints() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && d.Option.Dynamic != nil && d.Option.Dynamic.Uints {
		return true
	}
	return false
}

func (f *Field) DropdownDynamicFloats() bool {
	if d := f.View.Dropdown; d != nil && d.Option != nil && d.Option.Dynamic != nil && d.Option.Dynamic.Floats {
		return true
	}
	return false
}

func (f *Field) DropdownTranslateOptionMethod() string {
	sb := strings.Builder{}
	sb.WriteString("Translate")
	sb.WriteString(f.Name.Value)
	if f.View.Dropdown.Option.Dynamic != nil {
		sb.WriteString("Dynamic")
	}
	switch {
	case f.DropdownStrings(), f.DropdownDynamicStrings():
		sb.WriteString("Strings")
	case f.DropdownInts(), f.DropdownDynamicInts():
		sb.WriteString("Ints")
	case f.DropdownUints(), f.DropdownDynamicUints():
		sb.WriteString("Uints")
	case f.DropdownFloats(), f.DropdownDynamicFloats():
		sb.WriteString("Floats")
	}
	return sb.String()
}

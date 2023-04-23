package goalgenerator

import (
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"reflect"
	"strconv"
	"strings"
)

var DROPDOWN_KIND = [4]string{"strings", "ints", "uints", "floats"}

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
			HasOne *struct {
				Pkg,
				Name,
				Field string `yaml:",omitempty"`
			} `yaml:",omitempty"`
		} `yaml:",omitempty"`
		Calendar *struct {
			ShowTime,
			ShowIcon bool `yaml:",omitempty"`
		} `yaml:",omitempty"`
		Inline *struct {
			HasMany *struct {
				Pkg,
				Name string `yaml:",omitempty"`
			} `yaml:",omitempty"`
		} `yaml:",omitempty"`
		MultiSelect *struct {
			Many2Many *struct {
				Pkg,
				Name,
				Field string `yaml:",omitempty"`
			} `yaml:"many2many,omitempty"`
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
		if d := f.View.Dropdown; d.BelongTo != nil {
			if d.BelongTo.Pkg == "" {
				t = d.BelongTo.Name
			} else {
				t = strings.Join([]string{d.BelongTo.Pkg, d.BelongTo.Name}, ".")
			}
		} else if d.HasOne != nil {
			if d.HasOne.Pkg == "" {
				t = d.HasOne.Name
			} else {
				t = strings.Join([]string{d.HasOne.Pkg, d.HasOne.Name}, ".")
			}
		} else {
			for _, kind := range DROPDOWN_KIND {
				if f.Dropdown(kind, false) || f.Dropdown(kind, true) {
					t = f.kindToType(kind)
					break
				}
			}
		}
	case f.View.Inline != nil:
		if i := f.View.Inline; i.HasMany != nil {
			t = "[]"
			if i.HasMany.Pkg == "" {
				t += i.HasMany.Name
			} else {
				t += strings.Join([]string{i.HasMany.Pkg, i.HasMany.Name}, ".")
			}
		}
	case f.View.MultiSelect != nil:
		if ms := f.View.MultiSelect; ms.Many2Many != nil {
			t = "[]"
			if ms.Many2Many.Pkg == "" {
				t += ms.Many2Many.Name
			} else {
				t += strings.Join([]string{ms.Many2Many.Pkg, ms.Many2Many.Name}, ".")
			}
		}
	}
	return
}

func (f *Field) kindToType(kind string) (t string) {
	if kind == "floats" {
		t = "float64"
	} else {
		t = kind[:len(kind)-1]
	}
	return
}

func (f *Field) writeString(sb *strings.Builder, hasPrev *bool, value string, more ...string) {
	f.writeStringWithSep(sb, hasPrev, ",", value, more...)
}

func (f *Field) writeStringWithSep(sb *strings.Builder, hasPrev *bool, sep, value string, more ...string) {
	if *hasPrev {
		sb.WriteString(sep)
	} else {
		*hasPrev = true
	}
	sb.WriteString(value)
	for _, s := range more {
		sb.WriteString(s)
	}
}

func (f *Field) many2many(sb *strings.Builder, hasPrev *bool) {
	if ms := f.View.MultiSelect; ms != nil && ms.Many2Many != nil {
		f.writeStringWithSep(sb, hasPrev, ";", "many2many:", strings.ToLower(f.Model.Name.Value), "_", strings.ToLower(pluralize.NewClient().Plural(ms.Many2Many.Name)))
	}
}

func (f *Field) gorm(sb *strings.Builder) (primary bool, unique bool) {
	if f.Database != nil {
		primary = f.Database.PrimaryKey
		unique = f.Database.Unique
		sb.WriteString(`gorm:"`)
		var hasPrev bool
		if primary {
			f.writeStringWithSep(sb, &hasPrev, ";", "primarykey")
		}
		if unique {
			f.writeStringWithSep(sb, &hasPrev, ";", "unique")
		}
		if f.Database.Index {
			f.writeStringWithSep(sb, &hasPrev, ";", "index")
		}
		f.many2many(sb, &hasPrev)
		sb.WriteString(`" `)
	} else if ms := f.View.MultiSelect; ms != nil && ms.Many2Many != nil {
		sb.WriteString(`gorm:"`)
		var hasPrev bool
		f.many2many(sb, &hasPrev)
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
									switch vf.Name {
									case "Dropdown":
										switch cf.Name {
										case "BelongTo", "HasOne":
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
									case "Inline", "MultiSelect":
										p := e.FieldByName("Pkg").Interface().(string)
										if p == "" {
											p = m.Package.Value
										}
										n := e.FieldByName("Name").Interface()
										switch cf.Name {
										case "HasMany":
											f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name), "=", fmt.Sprintf("%s.%s", p, n))
										case "Many2Many":
											fn := e.FieldByName("Field").Interface()
											f.writeString(sb, &hasPrev, ToLowerFirstLetter(cf.Name), "=", fmt.Sprintf("%s.%s.%s", p, n, fn))
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

func (f *Field) Dropdown(kind string, dynamic bool) (ok bool) {
	if d := f.View.Dropdown; d != nil && d.Option != nil {
		if dynamic {
			if dyn := d.Option.Dynamic; dyn != nil {
				dfVal := reflect.ValueOf(dyn).Elem().FieldByName(ToUpperFirstLetter(kind))
				ok = dfVal.Bool()
			}
		} else {
			ofVal := reflect.ValueOf(d.Option).Elem().FieldByName(ToUpperFirstLetter(kind))
			ok = ofVal.Len() > 0
		}
	}
	return
}

func (f *Field) DropdownOptions(kind string) string {
	if f.Dropdown(kind, false) {
		t := f.kindToType(kind)
		sb := strings.Builder{}
		sb.WriteString("[]")
		sb.WriteString(t)
		sb.WriteRune('{')
		var hasPrev bool
		ofVal := reflect.ValueOf(f.View.Dropdown.Option).Elem().FieldByName(ToUpperFirstLetter(kind))
		for i := 0; i < ofVal.Len(); i++ {
			ofItemVal := ofVal.Index(i).FieldByName("Value")
			switch kind {
			case "strings":
				f.writeString(&sb, &hasPrev, `"`, ofItemVal.String(), `"`)
			case "ints":
				f.writeString(&sb, &hasPrev, fmt.Sprintf("%d", ofItemVal.Elem().Int()))
			case "uints":
				f.writeString(&sb, &hasPrev, fmt.Sprintf("%d", ofItemVal.Elem().Uint()))
			case "floats":
				f.writeString(&sb, &hasPrev, strconv.FormatFloat(ofItemVal.Elem().Float(), 'f', -1, 64))
			}
		}
		sb.WriteRune('}')
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownLabels(kind string) string {
	if f.Dropdown(kind, false) {
		sb := strings.Builder{}
		sb.WriteString("map[string]map[string]string{")
		for _, lang := range []string{"en", "zh-CN"} {
			sb.WriteRune('"')
			sb.WriteString(lang)
			sb.WriteRune('"')
			sb.WriteString(": {")

			var hasPrev bool
			ofVal := reflect.ValueOf(f.View.Dropdown.Option).Elem().FieldByName(ToUpperFirstLetter(kind))
			for i := 0; i < ofVal.Len(); i++ {
				ofItemVal := ofVal.Index(i)
				val := ofItemVal.FieldByName("Value")
				langVal := ofItemVal.FieldByName(ToUpperFirstLetter(strings.ReplaceAll(lang, "-", "")))
				var option string
				switch kind {
				case "strings":
					option = val.String()
				case "ints":
					option = fmt.Sprintf("%d", val.Elem().Int())
				case "uints":
					option = fmt.Sprintf("%d", val.Elem().Uint())
				case "floats":
					option = strconv.FormatFloat(val.Elem().Float(), 'f', -1, 64)
				}
				if option != "" {
					f.writeString(&sb, &hasPrev, `"`, option, `"`, ": ", `"`, langVal.String(), `"`)
				}
			}

			sb.WriteString(`},`)
		}
		sb.WriteRune('}')
		return sb.String()
	} else {
		return ""
	}
}

func (f *Field) DropdownTranslateOptionMethod() string {
	sb := strings.Builder{}
	sb.WriteString("Translate")
	sb.WriteString(f.Name.Value)
	if f.View.Dropdown.Option.Dynamic != nil {
		sb.WriteString("Dynamic")
	}
	for _, kind := range DROPDOWN_KIND {
		if f.Dropdown(kind, false) || f.Dropdown(kind, true) {
			sb.WriteString(ToUpperFirstLetter(kind))
			break
		}
	}
	return sb.String()
}

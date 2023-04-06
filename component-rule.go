package goalgenerator

import (
	"fmt"
	"strings"
)

const (
	COMPONENT_RULE_LAZY          ComponentRule = "lazy"
	COMPONENT_RULE_SORTABLE      ComponentRule = "sortable"
	COMPONENT_RULE_DESC          ComponentRule = "desc"
	COMPONENT_RULE_GLOBAL_SEARCH ComponentRule = "globalSearch"
	COMPONENT_RULE_SECRET        ComponentRule = "secret"
	COMPONENT_RULE_HIDDEN        ComponentRule = "hidden"
	COMPONENT_RULE_READONLY      ComponentRule = "readonly"
	COMPONENT_RULE_POSTONLY      ComponentRule = "postonly"
	COMPONENT_RULE_AUTOWIRED     ComponentRule = "autowired"
	COMPONENT_RULE_FILTER        ComponentRule = "filter"
	COMPONENT_RULE_SHOW_TIME     ComponentRule = "showTime"
	COMPONENT_RULE_SHOW_ICON     ComponentRule = "showIcon"
	COMPONENT_RULE_BELONGTO      ComponentRule = "@belongTo"
	COMPONENT_RULE_UPLOADTO      ComponentRule = "@uploadTo"
)

var FIELD_COMPONENT_RULES = []ComponentRule{
	COMPONENT_RULE_SORTABLE,
	COMPONENT_RULE_DESC,
	COMPONENT_RULE_GLOBAL_SEARCH,
	COMPONENT_RULE_HIDDEN,
	COMPONENT_RULE_READONLY,
	COMPONENT_RULE_POSTONLY,
	COMPONENT_RULE_AUTOWIRED,
	COMPONENT_RULE_FILTER,
	COMPONENT_RULE_SHOW_TIME,
	COMPONENT_RULE_SHOW_ICON,
	COMPONENT_RULE_BELONGTO,
	COMPONENT_RULE_UPLOADTO,
}

var MODEL_COMPONENT_RULES = []ComponentRule{
	COMPONENT_RULE_LAZY,
}

type ComponentRule string

func (cr *ComponentRule) ValidModel() error {
	return validComponentRule(cr, MODEL_COMPONENT_RULES)
}

func (cr *ComponentRule) ValidField() error {
	return validComponentRule(cr, FIELD_COMPONENT_RULES)
}

// IsProp implements IProp
func (cr *ComponentRule) IsProp() bool {
	return strings.HasPrefix(string(*cr), "@")
}

// Prop implements IProp
func (cr *ComponentRule) Prop() (string, string) {
	if cr.IsProp() {
		kv := strings.Split(string(*cr), "=")
		if len(kv) == 2 {
			return kv[0], kv[1]
		}
	}
	return "", ""
}

var _ IProp = (*ComponentRule)(nil)

func validComponentRule(cr *ComponentRule, rules []ComponentRule) error {
	for _, r := range rules {
		if *cr == r || cr.IsProp() && strings.HasPrefix(string(*cr), string(r)) {
			return nil
		}
	}
	return fmt.Errorf("invalid component rule: %v", *cr)
}

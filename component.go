package goalgenerator

import (
	"fmt"
	"strings"
)

const (
	COMPONENT_NUMBER   Component = "number"
	COMPONENT_UUID     Component = "uuid"
	COMPONENT_SWITCH   Component = "switch"
	COMPONENT_TEXT     Component = "text"
	COMPONENT_PASSWORD Component = "password"
	COMPONENT_DROPDOWN Component = "dropdown"
	COMPONENT_CALENDAR Component = "calendar"
	COMPONENT_FILE     Component = "file"
)

var COMPONENTS = []Component{
	COMPONENT_NUMBER,
	COMPONENT_UUID,
	COMPONENT_SWITCH,
	COMPONENT_TEXT,
	COMPONENT_PASSWORD,
	COMPONENT_DROPDOWN,
	COMPONENT_CALENDAR,
	COMPONENT_FILE,
}

type Component string

// IsProp implements IProp
func (c *Component) IsProp() bool {
	return strings.HasPrefix(string(*c), "@")
}

// Prop implements IProp
func (c *Component) Prop() (string, string) {
	if c.IsProp() {
		kv := strings.Split(string(*c), "=")
		if len(kv) == 2 {
			return kv[0], kv[1]
		}
	}
	return "", ""
}

// Valid implements IValid
func (c *Component) Valid() error {
	for _, component := range COMPONENTS {
		if *c == component || c.IsProp() && strings.HasPrefix(string(*c), string(component)) {
			return nil
		}
	}
	return fmt.Errorf("invalid component: %v", c)
}

var _ IValid = (*Component)(nil)
var _ IProp = (*Component)(nil)

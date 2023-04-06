package goalgenerator

import (
	"fmt"
	"strings"
)

const (
	VALIDATE_RULE_REQUIRED ValidateRule = "required"
	VALIDATE_RULE_EMAIL    ValidateRule = "email"
	VALIDATE_RULE_ALPHANUM ValidateRule = "alphanum"
	VALIDATE_RULE_ALPHA    ValidateRule = "alpha"
	VALIDATE_RULE_MIN      ValidateRule = "@min"
	VALIDATE_RULE_MAX      ValidateRule = "@max"
	VALIDATE_RULE_LEN      ValidateRule = "@len"
)

var VALIDATE_RULES = []ValidateRule{
	VALIDATE_RULE_REQUIRED,
	VALIDATE_RULE_EMAIL,
	VALIDATE_RULE_ALPHANUM,
	VALIDATE_RULE_ALPHA,
	VALIDATE_RULE_MIN,
	VALIDATE_RULE_MAX,
	VALIDATE_RULE_LEN,
}

type ValidateRule string

// IsProp implements IProp
func (vr *ValidateRule) IsProp() bool {
	return strings.HasPrefix(string(*vr), "@")
}

// Prop implements IProp
func (vr *ValidateRule) Prop() (string, string) {
	if vr.IsProp() {
		kv := strings.Split(string(*vr), "=")
		if len(kv) == 2 {
			return kv[0], kv[1]
		}
	}
	return "", ""
}

// Valid implements IValid
func (vr *ValidateRule) Valid() error {
	for _, r := range VALIDATE_RULES {
		if *vr == r || vr.IsProp() && strings.HasPrefix(string(*vr), string(r)) {
			return nil
		}
	}
	return fmt.Errorf("invalid validate rule: %v", vr)
}

var _ IValid = (*ValidateRule)(nil)
var _ IProp = (*ValidateRule)(nil)

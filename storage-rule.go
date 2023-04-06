package goalgenerator

import (
	"fmt"
	"strings"
)

const (
	STORAGE_RULE_PRIMARY        StorageRule = "primary"
	STORAGE_RULE_UNIQUE         StorageRule = "unique"
	STORAGE_RULE_INDEX          StorageRule = "index"
	STORAGE_RULE_EMBEDDING_BASE StorageRule = "embeddingBase"
)

var FIELD_STORAGE_RULES = []StorageRule{
	STORAGE_RULE_PRIMARY,
	STORAGE_RULE_UNIQUE,
	STORAGE_RULE_INDEX,
}

var MODEL_STORAGE_RULES = []StorageRule{
	STORAGE_RULE_EMBEDDING_BASE,
}

type StorageRule string

func (sr *StorageRule) ValidField() error {
	return validStorageRule(sr, FIELD_STORAGE_RULES)
}

func (sr *StorageRule) ValidModel() error {
	return validStorageRule(sr, MODEL_STORAGE_RULES)
}

// IsProp implements IProp
func (sr *StorageRule) IsProp() bool {
	return strings.HasPrefix(string(*sr), "@")
}

// Prop implements IProp
func (sr *StorageRule) Prop() (string, string) {
	if sr.IsProp() {
		kv := strings.Split(string(*sr), "=")
		if len(kv) == 2 {
			return kv[0], kv[1]
		}
	}
	return "", ""
}

var _ IProp = (*StorageRule)(nil)

func validStorageRule(sr *StorageRule, rules []StorageRule) error {
	for _, r := range rules {
		if *sr == r || sr.IsProp() && strings.HasPrefix(string(*sr), string(r)) {
			return nil
		}
	}
	return fmt.Errorf("invalid storage rule: %v", sr)
}

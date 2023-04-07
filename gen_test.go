package goalgenerator

import (
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOutput(t *testing.T) {
	m := &Model{
		Package: "auth",
		Name: "Role",
		StorageRules: []StorageRule{
			STORAGE_RULE_EMBEDDING_BASE,
		},
		ComponentRules: []ComponentRule{
			COMPONENT_RULE_LAZY,
		},
		Fields: []Field{
			{
				Name: "ID",
				StorageRules: []StorageRule{
					STORAGE_RULE_PRIMARY,
				},
				Component: COMPONENT_NUMBER,
				ComponentRules: []ComponentRule{
					COMPONENT_RULE_SORTABLE,
				},
			},
			{
				Name: "Name",
				StorageRules: []StorageRule{
					STORAGE_RULE_UNIQUE,
				},
				Component: COMPONENT_TEXT,
				ComponentRules: []ComponentRule{
					COMPONENT_RULE_SORTABLE,
				},
				ValidateRules: []ValidateRule{
					VALIDATE_RULE_REQUIRED,
					VALIDATE_RULE_ALPHANUM,
					ValidateRule("@min=3"),
					ValidateRule("@max=40"),
				},
			},
		},
	}

	if out, err := yaml.Marshal(m); err != nil {
		log.Fatal(err)
	} else {
		os.WriteFile("test/Role.yaml", out, os.ModePerm)
	}
}

func TestInput(t *testing.T) {
	in, _ := os.ReadFile("test/User.yaml")
	m := &Model{}
	if err := yaml.Unmarshal(in, m); err != nil {
		log.Fatal(err)
	} else {
		if m.Name != "User" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}

func TestGenRole(t *testing.T) {
	in, _ := os.ReadFile("test/Role.yaml")
	m := &Model{}
	if err := yaml.Unmarshal(in, m); err != nil {
		log.Fatal(err)
	} else {
		if m.Name != "Role" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := GenModel(m); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}

func TestGenUser(t *testing.T) {
	in, _ := os.ReadFile("test/User.yaml")
	m := &Model{}
	if err := yaml.Unmarshal(in, m); err != nil {
		log.Fatal(err)
	} else {
		if m.Name != "User" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := GenModel(m); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}

func TestGenSession(t *testing.T) {
	in, _ := os.ReadFile("test/Session.yaml")
	m := &Model{}
	if err := yaml.Unmarshal(in, m); err != nil {
		log.Fatal(err)
	} else {
		if m.Name != "Session" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := GenModel(m); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}
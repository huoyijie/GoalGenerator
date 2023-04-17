package goalgenerator

import (
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInput(t *testing.T) {
	in, _ := os.ReadFile("test/User.yaml")
	m := &Model{}
	if err := yaml.Unmarshal(in, m); err != nil {
		log.Fatal(err)
	} else {
		if m.Name.Value != "User" {
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
		if m.Name.Value != "Role" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := m.Gen(); err != nil {
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
		if m.Name.Value != "User" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := m.Gen(); err != nil {
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
		if m.Name.Value != "Session" {
			t.Fatal("Unmarshal failed")
		}
		if err := m.Valid(); err != nil {
			t.Fatalf("%+v", err)
		}
		if err := m.Gen(); err != nil {
			t.Fatalf("%+v", err)
		}
	}
}

package main

import (
	"flag"
	"path/filepath"

	goalgenerator "github.com/huoyijie/GoalGenerator"
	"gopkg.in/yaml.v3"

	"log"
	"os"
	"strings"
)

func main() {
	var yamlDir string
	flag.StringVar(&yamlDir, "d", "./", "Yaml files dictionaries")
	flag.Parse()

	genForDir(yamlDir)
}

func genForDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		name := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			genForDir(name)
		} else if strings.HasSuffix(entry.Name(), ".yaml") {
			in, _ := os.ReadFile(name)
			m := &goalgenerator.Model{}
			if err := yaml.Unmarshal(in, m); err != nil {
				log.Fatal(err)
			} else {
				if err := m.Valid(); err != nil {
					log.Fatal(err)
				}
				if err := goalgenerator.GenModel(m); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

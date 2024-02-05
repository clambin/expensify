package categorizer

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Rules map[Category][]Rule

func LoadRules(r io.Reader) (Rules, error) {
	var rules Rules
	err := yaml.NewDecoder(r).Decode(&rules)
	return rules, err
}

func LoadRulesFromFile(filename string) (Rules, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Rules{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	return LoadRules(file)
}

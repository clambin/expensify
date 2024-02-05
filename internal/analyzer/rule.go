package analyzer

import (
	"errors"
	"github.com/clambin/expensify/internal/payment"
	"github.com/clambin/expensify/pkg/maps"
	"gopkg.in/yaml.v3"
	"io"
	"regexp"
)

type Rule struct {
	Description string
}

func (r Rule) Match(p payment.Payment) bool {
	if r.Description != "" {
		if !p.Match(regexp.MustCompile(r.Description)) {
			return false
		}
	}
	return true
}

type Rules map[string][]Rule

func LoadRules(r io.Reader) (Rules, error) {
	var rules Rules
	err := yaml.NewDecoder(r).Decode(&rules)
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return rules, err
}

func (r Rules) Match(p payment.Payment) (string, bool) {
	for _, name := range maps.SortedKeys(r) {
		for _, _r := range r[name] {
			if _r.Match(p) {
				return name, true
			}
		}
	}
	return "", false
}

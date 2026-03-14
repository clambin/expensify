package rules

type Rule struct {
	Field    string   `yaml:"field,omitempty" json:"field,omitempty"`
	Operator Operator `yaml:"operator,omitempty" json:"operator,omitempty"`
	Value    any      `yaml:"value,omitempty" json:"value,omitempty"`
	Any      []Rule   `yaml:"any,omitempty" json:"any,omitempty"`
	All      []Rule   `yaml:"all,omitempty" json:"all,omitempty"`
}

func (r Rule) Evaluate(state map[string]any) (bool, error) {
	var anyRules, allRules, singleRule bool
	if len(r.Any) > 0 {
		for _, sub := range r.Any {
			ok, err := sub.Evaluate(state)
			if err != nil {
				return false, err
			}
			if ok {
				anyRules = true
				break
			}
		}
	}

	if len(r.All) > 0 {
		allRules = true
		for _, sub := range r.All {
			ok, err := sub.Evaluate(state)
			if err != nil {
				return false, err
			}
			if !ok {
				allRules = false
				break
			}
		}
	}

	if r.Operator != "" {
		if actual, ok := state[r.Field]; ok {
			var err error
			singleRule, err = r.Operator.Evaluate(actual, r.Value)
			if err != nil {
				return false, err
			}
		}
	}

	return anyRules || allRules || singleRule, nil
}

package hue

import (
	"encoding/json"
	"strconv"
)

type RuleActionsBody struct {
	Scene string `json:"scene"`
}

type RuleConditions struct {
	Address  string `json:"address"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type RuleActions struct {
	Address string          `json:"address"`
	Method  string          `json:"method"`
	Body    RuleActionsBody `json:"body"`
}

// Rule contains all data returned from the Phillips Hue API
// for an individual Phillips Hue rule
type Rule struct {
	Name           string           `json:"name"`
	LastTriggered  string           `json:"lasttriggered"`
	CreationTime   string           `json:"creationtime"`
	TimesTriggered int              `json:"timestriggered"`
	Owner          string           `json:"owner"`
	Status         string           `json:"status"`
	Conditions     []RuleConditions `json:"conditions"`
	Actions        []RuleActions    `json:"actions"`
	ID             int              `json:"id"`
}

// GetRules gets all Phillips Hue rules
func (h *Connection) GetRules() ([]Rule, error) {
	data, err := h.get("rules")
	if err != nil {
		return []Rule{}, err
	}

	if len(data) == 0 {
		return []Rule{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return []Rule{}, err
	}

	allRules := []Rule{}

	// Loop through all keys in the map and Unmarshal into
	// Rule type
	for key, val := range fullResponse {
		rule := Rule{}

		s, err := json.Marshal(val)
		if err != nil {
			return []Rule{}, err
		}

		err = json.Unmarshal(s, &rule)
		if err != nil {
			return []Rule{}, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return []Rule{}, err
		}

		rule.ID = id

		allRules = append(allRules, rule)
	}

	return allRules, nil
}

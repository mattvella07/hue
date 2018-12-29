package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// RuleActionsBody contains the data for the Body field in the
// RuleActions type
type RuleActionsBody struct {
	Scene string `json:"scene"`
}

// RuleConditions contains the data for the Conditions field in the
// Rule type
type RuleConditions struct {
	Address  string `json:"address"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// RuleActions contains the data for the Actions field in the
// Rule type
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

// GetRule gets the specified Phillips Hue rule
func (h *Connection) GetRule(rule int) (Rule, error) {
	data, err := h.get(fmt.Sprintf("rules/%d", rule))
	if err != nil {
		return Rule{}, err
	}

	// Rule not found
	if len(data) == 0 {
		return Rule{}, fmt.Errorf("Rule %d not found", rule)
	}

	ruleRes := Rule{}

	err = json.Unmarshal(data, &ruleRes)
	if err != nil {
		return Rule{}, err
	}

	return ruleRes, nil
}

// CreateRule creates a new rule with the specified name
func (h *Connection) CreateRule(name string, conditions []RuleConditions, actions []RuleActions) error {
	// Error checking
	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	bodyStr := fmt.Sprintf("{\"name\": \"%s\"", name)

	if len(conditions) > 0 && (strings.Trim(conditions[0].Address, " ") != "" || strings.Trim(conditions[0].Operator, " ") != "" || strings.Trim(conditions[0].Value, " ") != "") {
		bodyStr += fmt.Sprintf(", \"conditions\": %s", h.formatStruct(conditions))
	}

	if len(actions) > 0 && (strings.Trim(actions[0].Address, " ") != "" || strings.Trim(actions[0].Method, " ") != "") {
		bodyStr += fmt.Sprintf(", \"actions\": %s", h.formatStruct(actions))
	}

	bodyStr += "}"

	reqBody := strings.NewReader(bodyStr)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rules", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// RenameRule renames the specified Phillips Hue rule
func (h *Connection) RenameRule(rule int, name string) error {
	// Error checking
	if !h.doesRuleExist(rule) {
		return fmt.Errorf("Rule %d not found", rule)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	reqBody := strings.NewReader(fmt.Sprintf("{ \"name\": \"%s\" }", name))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rules/%d", h.baseURL, rule), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRule deletes a Phillips Hue rule from the bridge
func (h *Connection) DeleteRule(rule int) error {
	// Error checking
	if !h.doesRuleExist(rule) {
		return fmt.Errorf("Rule %d not found", rule)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/rules/%d", h.baseURL, rule), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) doesRuleExist(rule int) bool {
	// If GetRule returns an error, then the rule doesn't exist
	_, err := h.GetRule(rule)
	if err != nil {
		return false
	}

	return true
}

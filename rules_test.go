package hue

import (
	"testing"
)

func TestGetRules(t *testing.T) {
	t.Run("One rule found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		rules, err := h.GetRules()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(rules) != expected {
				t.Fatalf("Expected %d schedule, got %d", expected, len(rules))
			}
		}

		{
			expected := "Rule 1"
			if rules[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, rules[0].Name)
			}
		}

		{
			expected := 10
			if rules[0].TimesTriggered != expected {
				t.Fatalf("Expected TimesTriggered to equal %d, got %d", expected, rules[0].TimesTriggered)
			}
		}

		{
			expected := "abc"
			if rules[0].Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, rules[0].Owner)
			}
		}

		{
			expected := "enabled"
			if rules[0].Status != expected {
				t.Fatalf("Expected Status to equal %s, got %s", expected, rules[0].Status)
			}
		}
	})

	t.Run("No rules found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		rules, err := h.GetRules()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(rules) != expected {
				t.Fatalf("Expected %d rules, got %d", expected, len(rules))
			}
		}
	})
}

func TestGetRule(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		rule, err := h.GetRule(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Rule 1"
			if rule.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, rule.Name)
			}
		}

		{
			expected := 10
			if rule.TimesTriggered != expected {
				t.Fatalf("Expected TimesTriggered to equal %d, got %d", expected, rule.TimesTriggered)
			}
		}

		{
			expected := "abc"
			if rule.Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, rule.Owner)
			}
		}
	})

	t.Run("Not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetRule(1)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Rule 1 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestCreateRule(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	conditions := []RuleConditions{}

	actions := []RuleActions{}

	t.Run("Successful rule creation", func(t *testing.T) {
		err := h.CreateRule("new rule", conditions, actions)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.CreateRule("", conditions, actions)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Name must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameRule(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameRule(1, "Rule Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Rule doesn't exist", func(t *testing.T) {
		err := h.RenameRule(3, "Rule Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Rule 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.RenameRule(1, "")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Name must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteRule(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteRule(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Rule doesn't exist", func(t *testing.T) {
		err := h.DeleteRule(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Rule 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

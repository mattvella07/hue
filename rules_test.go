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

package hue

import (
	"testing"
)

type groupTestData struct {
	One Group `json:"1"`
}

func TestGetAllGroups(t *testing.T) {
	t.Run("One group found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		groups, err := h.GetAllGroups()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(groups) != expected {
				t.Fatalf("Expected %d group, got %d", expected, len(groups))
			}
		}

		{
			expected := "Group 1"
			if groups[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, groups[0].Name)
			}
		}

		{
			expected := 2
			if len(groups[0].Lights) != expected {
				t.Fatalf("Expected group to have %d lights, got %d", expected, len(groups[0].Lights))
			}
		}

		{
			expected := "LightGroup"
			if groups[0].Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, groups[0].Type)
			}
		}

		{
			expected := "ct"
			if groups[0].Action.ColorMode != expected {
				t.Fatalf("Expected ColorMode to equal %s, got %s", expected, groups[0].Action.ColorMode)
			}
		}
	})

	t.Run("No groups found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		groups, err := h.GetAllGroups()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(groups) != expected {
				t.Fatalf("Expected %d lights, got %d", expected, len(groups))
			}
		}
	})
}

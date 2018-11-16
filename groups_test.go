package hue

import (
	"testing"
)

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

func TestCreateGroup(t *testing.T) {
	t.Run("Successful group creation - LightGroup", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "LightGroup", "", []string{"1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Room", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "Room", "", []string{"1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Luminaire", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "Luminaire", "", []string{"1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - LightSource", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "LightSource", "", []string{"1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Empty group name", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "", "", []string{"1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid group name", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("", "LightGroup", "", []string{"1"})
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Name must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid group type", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "InvalidGroupType", "", []string{"1"})
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Group Type must be one of the following: LightGroup, Room, Luminaire, LightSource"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid light id", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.CreateGroup("New Group", "LightGroup", "", []string{"3"})
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

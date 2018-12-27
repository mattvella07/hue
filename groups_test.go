package hue

import (
	"testing"
)

func TestGetGroups(t *testing.T) {
	t.Run("One group found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		groups, err := h.GetGroups()
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

		groups, err := h.GetGroups()
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
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful group creation - LightGroup", func(t *testing.T) {
		err := h.CreateGroup("New Group", "LightGroup", "", []int{1})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Room", func(t *testing.T) {
		err := h.CreateGroup("New Group", "Room", "", []int{1})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Luminaire", func(t *testing.T) {
		err := h.CreateGroup("New Group", "Luminaire", "", []int{1})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - LightSource", func(t *testing.T) {
		err := h.CreateGroup("New Group", "LightSource", "", []int{1})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful group creation - Empty group name", func(t *testing.T) {
		err := h.CreateGroup("New Group", "", "", []int{1})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid group name", func(t *testing.T) {
		err := h.CreateGroup("", "LightGroup", "", []int{1})
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

	t.Run("Invalid group type", func(t *testing.T) {
		err := h.CreateGroup("New Group", "InvalidGroupType", "", []int{1})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group Type must be one of the following: LightGroup, Room, Luminaire, LightSource"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid light id", func(t *testing.T) {
		err := h.CreateGroup("New Group", "LightGroup", "", []int{3})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "One of the lights is invalid"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestGetGroup(t *testing.T) {
	t.Run("Group found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		group, err := h.GetGroup(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Group 1"
			if group.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, group.Name)
			}
		}

		{
			expected := "LightGroup"
			if group.Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, group.Type)
			}
		}
	})

	t.Run("Group not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetGroup(2)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 2 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameGroup(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameGroup(1, "Group Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.RenameGroup(3, "Group Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.RenameGroup(1, "")
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

func TestSetLightsInGroup(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.SetLightsInGroup(1, []int{2})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.SetLightsInGroup(3, []int{2})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid light", func(t *testing.T) {
		err := h.SetLightsInGroup(1, []int{6})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "One of the lights is invalid"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestSetGroupClass(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.SetGroupClass(1, "Other")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.SetGroupClass(3, "Other")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid class", func(t *testing.T) {
		err := h.SetGroupClass(1, "")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Class must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteGroup(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteGroup(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.DeleteGroup(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOnGroup(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.TurnOnGroup(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.TurnOnGroup(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOnGroupWithColor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 0.3, 0.2, 100, 200, 233)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(3, 0.3, 0.2, 100, 200, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid x value", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 2, 0.2, 100, 200, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Invalid color value: x must be between 0 and 1"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid y value", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 0.2, 3, 100, 200, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Invalid color value: y must be between 0 and 1"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid bri value", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 0.3, 0.2, 300, 200, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Invalid brightness value: bri must be between 1 and 254"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid hue value", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 0.3, 0.2, 100, 65539, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Invalid hue value: hue must be between 0 and 65,535"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid sat value", func(t *testing.T) {
		err := h.TurnOnGroupWithColor(1, 0.3, 0.2, 100, 200, 350)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Invalid saturation value: sat must be between 0 and 254"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOffGroup(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.TurnOffGroup(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Group doesn't exist", func(t *testing.T) {
		err := h.TurnOffGroup(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Group 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

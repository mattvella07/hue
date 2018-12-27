package hue

import (
	"testing"
)

func TestGetLights(t *testing.T) {
	t.Run("One light found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		lights, err := h.GetLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(lights) != expected {
				t.Fatalf("Expected %d light, got %d", expected, len(lights))
			}
		}

		{
			expected := "Hue color lamp 1"
			if lights[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, lights[0].Name)
			}
		}

		{
			expected := "Extended color"
			if lights[0].Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, lights[0].Type)
			}
		}

		{
			expected := "Phillips"
			if lights[0].ManufacturerName != expected {
				t.Fatalf("Expected ManufacturerName to equal %s, got %s", expected, lights[0].ManufacturerName)
			}
		}
	})

	t.Run("No lights found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		lights, err := h.GetLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(lights) != expected {
				t.Fatalf("Expected %d lights, got %d", expected, len(lights))
			}
		}
	})
}

func TestGetNewLights(t *testing.T) {
	t.Run("New light found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		newLights, err := h.GetNewLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(newLights.NewLights) != expected {
				t.Fatalf("Expected number of new lights to equal %d, got %d", expected, len(newLights.NewLights))
			}
		}

		{
			expected := 5
			if newLights.NewLights[0].ID != expected {
				t.Fatalf("Expected ID of first new light to equal %d, got %d", expected, newLights.NewLights[0].ID)
			}
		}

		{
			expected := "Hue lamp 5"
			if newLights.NewLights[0].Name != expected {
				t.Fatalf("Expected Name of first new light to equal %s, got %s", expected, newLights.NewLights[0].Name)
			}
		}

		{
			expected := "2018-10-12T12:00:00"
			if newLights.LastScan != expected {
				t.Fatalf("Expected LastScan to equal %s, got %s", expected, newLights.LastScan)
			}
		}
	})

	t.Run("No new lights found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		newLights, err := h.GetNewLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(newLights.NewLights) != expected {
				t.Fatalf("Expected %d new lights, got %d", expected, len(newLights.NewLights))
			}
		}
	})
}

func TestFindNewLights(t *testing.T) {
	h, server := createTestConnection(3)
	defer server.Close()

	err := h.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLight(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		light, err := h.GetLight(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Hue color lamp 1"
			if light.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, light.Name)
			}
		}

		{
			expected := "Extended color"
			if light.Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, light.Type)
			}
		}

		{
			expected := "Phillips"
			if light.ManufacturerName != expected {
				t.Fatalf("Expected ManufacturerName to equal %s, got %s", expected, light.ManufacturerName)
			}
		}
	})

	t.Run("Not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetLight(1)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 1 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameLight(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameLight(1, "Light Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		err := h.RenameLight(3, "Light Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Inavlid name", func(t *testing.T) {
		err := h.RenameLight(1, "")
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

func TestTurnOnLight(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Light exists", func(t *testing.T) {
		err := h.TurnOnLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		err := h.TurnOnLight(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOnLightWithColor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Light exists", func(t *testing.T) {
		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 200, 233)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		err := h.TurnOnLightWithColor(3, 0.3, 0.2, 100, 200, 233)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid x value", func(t *testing.T) {
		err := h.TurnOnLightWithColor(1, 2, 0.2, 100, 200, 233)
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
		err := h.TurnOnLightWithColor(1, 0.2, 3, 100, 200, 233)
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
		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 300, 200, 233)
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
		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 65539, 233)
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
		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 200, 350)
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

func TestTurnOffLight(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Light exists", func(t *testing.T) {
		err := h.TurnOffLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		err := h.TurnOffLight(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteLight(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful delete", func(t *testing.T) {
		err := h.DeleteLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		err := h.DeleteLight(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Light 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

package hue

import "testing"

func TestCreateUser(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful user creation", func(t *testing.T) {
		err := h.CreateUser("app#device")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid deviceType", func(t *testing.T) {
		err := h.CreateUser("")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "deviceType must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestGetConfiguration(t *testing.T) {
	t.Run("One group found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		config, err := h.GetConfiguration()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Phillips hue"
			if config.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, config.Name)
			}
		}

		{
			expected := "abcd:efgh"
			if config.Mac != expected {
				t.Fatalf("Expected Mac to equal %s, got %s", expected, config.Mac)
			}
		}

		{
			expected := 15
			if config.ZigbeeChannel != expected {
				t.Fatalf("Expected ZigbeeChannel to equal %d, got %d", expected, config.ZigbeeChannel)
			}
		}

		{
			expected := "192.0.0.1"
			if config.IPAddress != expected {
				t.Fatalf("Expected IpAddress to equal %s, got %s", expected, config.IPAddress)
			}
		}
	})

	t.Run("No config found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		config, err := h.GetConfiguration()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := ""
			if config.Name != expected {
				t.Fatalf("Expected Name to be empty, got %s", config.Name)
			}
		}
	})
}

func TestDeleteUser(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteUser("1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("User empty", func(t *testing.T) {
		err := h.DeleteUser("")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "User must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

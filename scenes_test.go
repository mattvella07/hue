package hue

import "testing"

func TestGetAllScenes(t *testing.T) {
	t.Run("One scene found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		scenes, err := h.GetAllScenes()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(scenes) != expected {
				t.Fatalf("Expected %d scene, got %d", expected, len(scenes))
			}
		}

		{
			expected := "Night time"
			if scenes[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, scenes[0].Name)
			}
		}

		{
			expected := "LightScene"
			if scenes[0].Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, scenes[0].Type)
			}
		}

		{
			expected := "abcd"
			if scenes[0].Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, scenes[0].Owner)
			}
		}

		{
			expected := false
			if scenes[0].Locked != expected {
				t.Fatalf("Expected Locked to equal %t, got %t", expected, scenes[0].Locked)
			}
		}
	})

	t.Run("No scenes found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		scenes, err := h.GetAllScenes()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(scenes) != expected {
				t.Fatalf("Expected %d scenes, got %d", expected, len(scenes))
			}
		}
	})
}

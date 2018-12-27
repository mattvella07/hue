package hue

import "testing"

func TestGetScenes(t *testing.T) {
	t.Run("One scene found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		scenes, err := h.GetScenes()
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

		scenes, err := h.GetScenes()
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

func TestCreateLightScene(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful scene creation", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateLightScene("New Scene", []int{1, 2}, true, appData)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful scene creation - empty Name", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateLightScene("", []int{1, 2}, true, appData)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Empty light slice", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateLightScene("New Scene", []int{}, true, appData)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Lights must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid light", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateLightScene("New Scene", []int{3}, true, appData)
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

func TestCreateGroupScene(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful scene creation", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateGroupScene("New Scene", 1, true, appData)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful scene creation - empty Name", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateGroupScene("", 1, true, appData)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid group", func(t *testing.T) {
		appData := SceneAppData{
			Version: 1,
			Data:    "data",
		}

		err := h.CreateGroupScene("New Scene", 3, true, appData)
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

func TestRenameScene(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameScene("1", "Scene Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Scene doesn't exist", func(t *testing.T) {
		err := h.RenameScene("3", "Scene Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Scene 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.RenameScene("1", "")
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

func TestSetLightsInScene(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.SetLightsInScene("1", []int{2})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Scene doesn't exist", func(t *testing.T) {
		err := h.SetLightsInScene("3", []int{2})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Scene 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid light", func(t *testing.T) {
		err := h.SetLightsInScene("1", []int{6})
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

func TestDeleteScene(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteScene("1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Scene doesn't exist", func(t *testing.T) {
		err := h.DeleteScene("3")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Scene 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestGetScene(t *testing.T) {
	t.Run("Scene found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		scene, err := h.GetScene("1")
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Night time"
			if scene.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, scene.Name)
			}
		}

		{
			expected := "LightScene"
			if scene.Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, scene.Type)
			}
		}

		{
			expected := "abcd"
			if scene.Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, scene.Owner)
			}
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetScene("2")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Scene 2 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

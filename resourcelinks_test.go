package hue

import "testing"

func TestGetResourceLinks(t *testing.T) {
	t.Run("One resource link found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		links, err := h.GetResourceLinks()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(links) != expected {
				t.Fatalf("Expected %d resource link, got %d", expected, len(links))
			}
		}

		{
			expected := "Sunrise"
			if links[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, links[0].Name)
			}
		}

		{
			expected := 1
			if links[0].Class != expected {
				t.Fatalf("Expected Class to equal %d, got %d", expected, links[0].Class)
			}
		}

		{
			expected := "abcd"
			if links[0].Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, links[0].Owner)
			}
		}

		{
			expected := "Wake up experience"
			if links[0].Description != expected {
				t.Fatalf("Expected Status to equal %s, got %s", expected, links[0].Description)
			}
		}
	})

	t.Run("No resource links found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		links, err := h.GetResourceLinks()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(links) != expected {
				t.Fatalf("Expected %d resource links, got %d", expected, len(links))
			}
		}
	})
}

func TestGetResourceLink(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		link, err := h.GetResourceLink(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Sunset"
			if link.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, link.Name)
			}
		}

		{
			expected := 1
			if link.Class != expected {
				t.Fatalf("Expected TimesTriggered to equal %d, got %d", expected, link.Class)
			}
		}

		{
			expected := "abcd"
			if link.Owner != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, link.Owner)
			}
		}

		{
			expected := "Go to sleep experience"
			if link.Description != expected {
				t.Fatalf("Expected Owner to equal %s, got %s", expected, link.Description)
			}
		}
	})

	t.Run("Not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetResourceLink(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Resource link 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestCreateResourceLink(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful resource link creation", func(t *testing.T) {
		err := h.CreateResourceLink("new resource link", "desc", true, []string{"/path/1"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.CreateResourceLink("", "desc", true, []string{"/path/1"})
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

	t.Run("Invalid links", func(t *testing.T) {
		err := h.CreateResourceLink("new resource link", "desc", true, []string{})
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Links must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameResourceLink(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameResourceLink(1, "Resource link Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Resource link doesn't exist", func(t *testing.T) {
		err := h.RenameResourceLink(3, "Resource link Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Resource link 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.RenameResourceLink(1, "")
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

func TestSetResourceLinkDescription(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Success", func(t *testing.T) {
		err := h.SetResourceLinkDescription(1, "New description")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Resource link doesn't exist", func(t *testing.T) {
		err := h.SetResourceLinkDescription(3, "New description")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Resource link 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid description", func(t *testing.T) {
		err := h.SetResourceLinkDescription(1, "")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Description must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteResourceLink(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteResourceLink(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Resource link doesn't exist", func(t *testing.T) {
		err := h.DeleteResourceLink(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Resource link 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

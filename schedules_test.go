package hue

import (
	"testing"
)

func TestGetAllSchedules(t *testing.T) {
	t.Run("One schedule found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		schedules, err := h.GetAllSchedules()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(schedules) != expected {
				t.Fatalf("Expected %d schedule, got %d", expected, len(schedules))
			}
		}

		{
			expected := "Timer"
			if schedules[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, schedules[0].Name)
			}
		}

		{
			expected := "Simple timer"
			if schedules[0].Description != expected {
				t.Fatalf("Expected Description to equal %s, got %s", expected, schedules[0].Description)
			}
		}

		{
			expected := "enabled"
			if schedules[0].Status != expected {
				t.Fatalf("Expected Status to equal %s, got %s", expected, schedules[0].Status)
			}
		}

		{
			expected := false
			if schedules[0].AutoDelete != expected {
				t.Fatalf("Expected AutoDelete to equal %t, got %t", expected, schedules[0].AutoDelete)
			}
		}
	})

	t.Run("No schedules found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		schedules, err := h.GetAllSchedules()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(schedules) != expected {
				t.Fatalf("Expected %d schedules, got %d", expected, len(schedules))
			}
		}
	})
}

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

func TestCreateSchedule(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful schedule creation", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("new schedule", "a new schedule", cmd, "2018-01-01", "enabled", true, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful schedule creation - empty name", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("", "a new schedule", cmd, "2018-01-01", "enabled", true, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Successful schedule creation - empty description", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("new schedule", "", cmd, "2018-01-01", "enabled", true, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid command address", func(t *testing.T) {
		cmd := ScheduleCommand{}

		err := h.CreateSchedule("new schedule", "a new schedule", cmd, "2018-01-01", "enabled", true, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Command Address must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid command method", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "GET",
		}

		err := h.CreateSchedule("new schedule", "a new schedule", cmd, "2018-01-01", "enabled", true, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Command Method must be either POST, PUT, or DELETE"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid command body", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("new schedule", "a new schedule", cmd, "2018-01-01", "enabled", true, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Command Body must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid localtime", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("new schedule", "description", cmd, "", "enabled", true, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Localtime must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid status", func(t *testing.T) {
		cmd := ScheduleCommand{
			Address: "abc",
			Body: ScheduleCommandBody{
				Scene: "123",
			},
			Method: "PUT",
		}

		err := h.CreateSchedule("new schedule", "description", cmd, "2018-01-01", "invalid", true, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Status must be either enabled or disabled"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestGetSchedule(t *testing.T) {
	t.Run("Schedule found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		schedule, err := h.GetSchedule(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Timer"
			if schedule.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, schedule.Name)
			}
		}

		{
			expected := "Simple timer"
			if schedule.Description != expected {
				t.Fatalf("Expected Description to equal %s, got %s", expected, schedule.Description)
			}
		}

		{
			expected := "enabled"
			if schedule.Status != expected {
				t.Fatalf("Expected Status to equal %s, got %s", expected, schedule.Status)
			}
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetSchedule(2)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Schedule not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteSchedule(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteSchedule(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Schedule doesn't exist", func(t *testing.T) {
		err := h.DeleteSchedule(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Schedule 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

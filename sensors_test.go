package hue

import (
	"testing"
)

func TestGetSensors(t *testing.T) {
	t.Run("One sensor found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		sensors, err := h.GetSensors()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(sensors) != expected {
				t.Fatalf("Expected %d schedule, got %d", expected, len(sensors))
			}
		}

		{
			expected := "Daylight"
			if sensors[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, sensors[0].Name)
			}
		}

		{
			expected := "PHDL00"
			if sensors[0].ModelID != expected {
				t.Fatalf("Expected ModelID to equal %s, got %s", expected, sensors[0].ModelID)
			}
		}

		{
			expected := "1.0"
			if sensors[0].SWVersion != expected {
				t.Fatalf("Expected SWVersion to equal %s, got %s", expected, sensors[0].SWVersion)
			}
		}

		{
			expected := "Phillips"
			if sensors[0].ManufacturerName != expected {
				t.Fatalf("Expected ManufacturerName to equal %s, got %s", expected, sensors[0].ManufacturerName)
			}
		}
	})

	t.Run("No sensors found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		sensors, err := h.GetSensors()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(sensors) != expected {
				t.Fatalf("Expected %d sensors, got %d", expected, len(sensors))
			}
		}
	})
}

func TestCreateSensor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	state := SensorState{
		Daylight:    true,
		LastUpdated: "2014-06-27T07:38:51",
	}

	config := SensorConfig{
		On:            true,
		Long:          "none",
		Lat:           "none",
		SunriseOffset: 20,
		SunsetOffset:  20,
	}

	t.Run("Successful sensor creation", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "SENSOR1", "1.0", "S", "abcd", "Phillips", state, config, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.CreateSensor("", "SENSOR1", "1.0", "S", "abcd", "Phillips", state, config, true)
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

	t.Run("Invalid modelid", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "", "1.0", "S", "abcd", "Phillips", state, config, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "ModelID must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid swversion", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "SENSOR1", "", "S", "abcd", "Phillips", state, config, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "SWVersion must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid type", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "SENSOR1", "1.0", "", "abcd", "Phillips", state, config, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Type must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid uniqueid", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "SENSOR1", "1.0", "S", "", "Phillips", state, config, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "UniqueID must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid manufacturername", func(t *testing.T) {
		err := h.CreateSensor("new sensor", "SENSOR1", "1.0", "S", "abcd", "", state, config, true)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "ManufacturerName must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestFindNewSensors(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	err := h.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetNewSensors(t *testing.T) {
	t.Run("New sensor found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		newSensors, err := h.GetNewSensors()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(newSensors.NewSensors) != expected {
				t.Fatalf("Expected number of new sensors to equal %d, got %d", expected, len(newSensors.NewSensors))
			}
		}

		{
			expected := 5
			if newSensors.NewSensors[0].ID != expected {
				t.Fatalf("Expected ID of first new light to equal %d, got %d", expected, newSensors.NewSensors[0].ID)
			}
		}

		{
			expected := "Sensor 5"
			if newSensors.NewSensors[0].Name != expected {
				t.Fatalf("Expected Name of first new light to equal %s, got %s", expected, newSensors.NewSensors[0].Name)
			}
		}

		{
			expected := "2018-10-12T12:00:00"
			if newSensors.LastScan != expected {
				t.Fatalf("Expected LastScan to equal %s, got %s", expected, newSensors.LastScan)
			}
		}
	})

	t.Run("No new sensors found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		newSensors, err := h.GetNewSensors()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(newSensors.NewSensors) != expected {
				t.Fatalf("Expected %d new lights, got %d", expected, len(newSensors.NewSensors))
			}
		}
	})
}

func TestGetSensor(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		sensor, err := h.GetSensor(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Daylight"
			if sensor.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, sensor.Name)
			}
		}

		{
			expected := "PHDL00"
			if sensor.ModelID != expected {
				t.Fatalf("Expected ModelID to equal %s, got %s", expected, sensor.ModelID)
			}
		}

		{
			expected := "1.0"
			if sensor.SWVersion != expected {
				t.Fatalf("Expected SWVersion to equal %s, got %s", expected, sensor.SWVersion)
			}
		}
	})

	t.Run("Not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetSensor(1)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Sensor 1 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameSensor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful rename", func(t *testing.T) {
		err := h.RenameSensor(1, "Sensor Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Sensor doesn't exist", func(t *testing.T) {
		err := h.RenameSensor(3, "Sensor Renamed")
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Sensor 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		err := h.RenameSensor(1, "")
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

func TestDeleteSensor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Successful", func(t *testing.T) {
		err := h.DeleteSensor(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Sensor doesn't exist", func(t *testing.T) {
		err := h.DeleteSensor(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Sensor 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOnSensor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Sensor exists", func(t *testing.T) {
		err := h.TurnOnSensor(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Sensor doesn't exist", func(t *testing.T) {
		err := h.TurnOnSensor(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Sensor 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOffSensor(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	t.Run("Sensor exists", func(t *testing.T) {
		err := h.TurnOffSensor(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Sensor doesn't exist", func(t *testing.T) {
		err := h.TurnOffSensor(3)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}

		{
			expected := "Sensor 3 not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

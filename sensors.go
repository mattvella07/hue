package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// SensorState contains the data for the State field in the
// Sensor type
type SensorState struct {
	Daylight    bool   `json:"daylight"`
	LastUpdated string `json:"lastupdated"`
}

// SensorConfig contains the data for the Config field in the
// Sensor type
type SensorConfig struct {
	On            bool   `json:"on"`
	Long          string `json:"long"`
	Lat           string `json:"lat"`
	SunriseOffset int    `json:"sunriseoffset"`
	SunsetOffset  int    `json:"sunsetoffset"`
}

// Sensor contains all data returned from the Phillips Hue API
// for an individual Phillips Hue sensor
type Sensor struct {
	State            SensorState  `json:"state"`
	Config           SensorConfig `json:"config"`
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	ModelID          string       `json:"modelid"`
	ManufacturerName string       `json:"manufacturername"`
	SWVersion        string       `json:"swversion"`
	ID               int          `json:"id"`
}

// NewSensor contains all data for a new Phillips Hue sensor
type NewSensor struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// NewSensorResponse is the response from the Phillips Hue API for
// all new sensors
type NewSensorResponse struct {
	NewSensors []NewSensor `json:"newSensors"`
	LastScan   string      `json:"lastscan"`
}

// GetSensors gets all Phillips Hue sensors
func (h *Connection) GetSensors() ([]Sensor, error) {
	data, err := h.get("sensors")
	if err != nil {
		return []Sensor{}, err
	}

	if len(data) == 0 {
		return []Sensor{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return []Sensor{}, err
	}

	allSensors := []Sensor{}

	// Loop through all keys in the map and Unmarshal into
	// Sensor type
	for key, val := range fullResponse {
		sensor := Sensor{}

		s, err := json.Marshal(val)
		if err != nil {
			return []Sensor{}, err
		}

		err = json.Unmarshal(s, &sensor)
		if err != nil {
			return []Sensor{}, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return []Sensor{}, err
		}

		sensor.ID = id

		allSensors = append(allSensors, sensor)
	}

	return allSensors, nil
}

// CreateSensor creates a new sensor with the specified name
func (h *Connection) CreateSensor(name, modelID, swVersion, sensorType, uniqueID, manufacturerName string, state SensorState, config SensorConfig, recycle bool) error {
	// Error checking
	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	if strings.Trim(modelID, " ") == "" {
		return errors.New("ModelID must not be empty")
	}

	if strings.Trim(swVersion, " ") == "" {
		return errors.New("SWVersion must not be empty")
	}

	if strings.Trim(sensorType, " ") == "" {
		return errors.New("Type must not be empty")
	}

	if strings.Trim(uniqueID, " ") == "" {
		return errors.New("UniqueID must not be empty")
	}

	if strings.Trim(manufacturerName, " ") == "" {
		return errors.New("ManufacturerName must not be empty")
	}

	reqBody := strings.NewReader(fmt.Sprintf("{\"name\": \"%s\", \"modelid\": \"%s\", \"swversion\": %s, \"type\": \"%s\", \"uniqueid\": \"%s\", \"manufacturername\": \"%s\", \"state\": %s, \"config\": %s, \"recycle\": %t }", name, modelID, swVersion, sensorType, uniqueID, manufacturerName, h.formatStruct(state), h.formatStruct(config), recycle))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sensors", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// FindNewSensors finds new Phillips Hue sensors that have been added since
// the last time performing this call
func (h *Connection) FindNewSensors() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sensors", h.baseURL), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// GetNewSensors gets Phillips Hue sensors that were discovered since the last time
// FindNewSensors was called
func (h *Connection) GetNewSensors() (NewSensorResponse, error) {
	data, err := h.get("sensors/new")
	if err != nil {
		return NewSensorResponse{}, err
	}

	if len(data) == 0 {
		return NewSensorResponse{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return NewSensorResponse{}, err
	}

	newSensorRes := NewSensorResponse{}

	// Loop through all keys in the map and Unmarshal into
	// Group type
	for key, val := range fullResponse {
		if key == "lastscan" {
			newSensorRes.LastScan = val.(string)
		} else {
			newSensor := NewSensor{}

			l, err := json.Marshal(val)
			if err != nil {
				return NewSensorResponse{}, err
			}

			err = json.Unmarshal(l, &newSensor)
			if err != nil {
				return NewSensorResponse{}, err
			}

			id, err := strconv.Atoi(key)
			if err != nil {
				return NewSensorResponse{}, err
			}

			newSensor.ID = id

			newSensorRes.NewSensors = append(newSensorRes.NewSensors, newSensor)
		}
	}

	return newSensorRes, nil
}

// GetSensor gets the specified Phillips Hue sensor
func (h *Connection) GetSensor(sensor int) (Sensor, error) {
	data, err := h.get(fmt.Sprintf("sensors/%d", sensor))
	if err != nil {
		return Sensor{}, err
	}

	// Sensor not found
	if len(data) == 0 {
		return Sensor{}, fmt.Errorf("Sensor %d not found", sensor)
	}

	sensorRes := Sensor{}

	err = json.Unmarshal(data, &sensorRes)
	if err != nil {
		return Sensor{}, err
	}

	return sensorRes, nil
}

// RenameSensor renames the specified Phillips Hue sensor
func (h *Connection) RenameSensor(sensor int, name string) error {
	// Error checking
	if !h.doesSensorExist(sensor) {
		return fmt.Errorf("Sensor %d not found", sensor)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	reqBody := strings.NewReader(fmt.Sprintf("{ \"name\": \"%s\" }", name))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/sensors/%d", h.baseURL, sensor), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSensor deletes a Phillips Hue sensor from the bridge
func (h *Connection) DeleteSensor(sensor int) error {
	// Error checking
	if !h.doesSensorExist(sensor) {
		return fmt.Errorf("Sensor %d not found", sensor)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/sensors/%d", h.baseURL, sensor), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// TurnOnSensor turns on the specified Phillips Hue sensor
func (h *Connection) TurnOnSensor(sensor int) error {
	// Error checking
	if !h.doesSensorExist(sensor) {
		return fmt.Errorf("Sensor %d not found", sensor)
	}

	reqBody := strings.NewReader("{ \"on\": true }")
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/sensors/%d/config", h.baseURL, sensor), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// TurnOffSensor turns off the specified Phillips Hue sensor
func (h *Connection) TurnOffSensor(sensor int) error {
	// Error checking
	if !h.doesSensorExist(sensor) {
		return fmt.Errorf("Sensor %d not found", sensor)
	}

	reqBody := strings.NewReader("{ \"on\": false }")
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/sensors/%d/config", h.baseURL, sensor), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) doesSensorExist(sensor int) bool {
	// If GetSensor returns an error, then the sensor doesn't exist
	_, err := h.GetSensor(sensor)
	if err != nil {
		return false
	}

	return true
}

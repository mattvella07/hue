package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ScheduleCommandBody contains the data for the Body field in
//  the ScheduleCommand type
type ScheduleCommandBody struct {
	Scene string `json:"scene"`
	Flag  bool   `json:"flag"`
	On    bool   `json:"on"`
}

// ScheduleCommand contains the data for the Command field in the
// Schedule type
type ScheduleCommand struct {
	Address string              `json:"address"`
	Body    ScheduleCommandBody `json:"body"`
	Method  string              `json:"method"`
}

// Schedule contains all data returned from the Phillips Hue API
// for an individual Phillips Hue schedule
type Schedule struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Command     ScheduleCommand `json:"command"`
	Time        string          `json:"time"`
	Created     string          `json:"created"`
	Status      string          `json:"status"`
	AutoDelete  bool            `json:"autodelete"`
	StartTime   string          `json:"starttime"`
	ID          int             `json:"id"`
}

// GetAllSchedules gets all Phillips Hue schedules
func (h *Connection) GetAllSchedules() ([]Schedule, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/schedules", h.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return []Schedule{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(body, &fullResponse)
	if err != nil {
		return nil, err
	}

	allSchedules := []Schedule{}

	// Loop through all keys in the map and Unmarshal into
	// Schedule type
	for key, val := range fullResponse {
		schedule := Schedule{}

		s, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(s, &schedule)
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}

		schedule.ID = id

		allSchedules = append(allSchedules, schedule)
	}

	return allSchedules, nil
}

// CreateSchedule creates a new schedule with the specified name
func (h *Connection) CreateSchedule(name, description string, command ScheduleCommand, localtime, status string, autodelete, recycle bool) error {
	// Error checking
	if &command == nil {
		return errors.New("Command must not be empty")
	}

	if strings.Trim(command.Address, " ") == "" {
		return errors.New("Command Address must not be empty")
	}

	if command.Method != "POST" && command.Method != "PUT" && command.Method != "DELETE" {
		return errors.New("Command Method must be either POST, PUT, or DELETE")
	}

	if strings.Trim(command.Body.Scene, " ") == "" {
		return errors.New("Command Body must not be empty")
	}

	if strings.Trim(localtime, " ") == "" {
		return errors.New("Localtime must not be empty")
	}

	if strings.Trim(status, " ") != "" {
		if status != "enabled" && status != "disabled" {
			return errors.New("Status must be either enabled or disabled")
		}
	}

	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	reqBody := strings.NewReader(fmt.Sprintf("{\"name\": \"%s\", \"description\": \"%s\", \"command\": %s, \"localtime\": \"%s\", \"status\": \"%s\", \"autodelete\": %t, \"recycle\": %t }", name, description, h.formatStruct(command), localtime, status, autodelete, recycle))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/schedules", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fullResponse := string(body)
	fullResponse = strings.ToLower(fullResponse)

	if strings.Contains(fullResponse, "error") {
		errMsg := fullResponse[strings.Index(fullResponse, "\"description\":\"")+15 : strings.Index(fullResponse, "\"}}")]
		return errors.New(errMsg)
	}

	return nil
}

// GetSchedule gets the specified all Phillips Hue schedule by ID
func (h *Connection) GetSchedule(schedule int) (Schedule, error) {
	err := h.initializeHue()
	if err != nil {
		return Schedule{}, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/schedules/%d", h.baseURL, schedule))
	if err != nil {
		return Schedule{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Schedule{}, err
	}

	// Schedule not found
	if len(body) == 0 {
		return Schedule{}, errors.New("Schedule not found")
	}

	scheduleRes := Schedule{}

	err = json.Unmarshal(body, &scheduleRes)
	if err != nil {
		return Schedule{}, err
	}

	return scheduleRes, nil
}

// RenameSchedule renames the specified Phillips Hue schedule
func (h *Connection) RenameSchedule(schedule int, name string) error {
	// Error checking
	if !h.doesScheduleExist(schedule) {
		return fmt.Errorf("Schedule %d not found", schedule)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	attributes := fmt.Sprintf("{ \"name\": \"%s\" }", name)

	err := h.updateSchedule(schedule, attributes)
	if err != nil {
		return err
	}

	return nil
}

// SetScheduleDescription sets the description for the specified Phillips Hue schedule
func (h *Connection) SetScheduleDescription(schedule int, description string) error {
	// Error checking
	if !h.doesScheduleExist(schedule) {
		return fmt.Errorf("Schedule %d not found", schedule)
	}

	attributes := fmt.Sprintf("{ \"description\": \"%s\" }", description)

	err := h.updateSchedule(schedule, attributes)
	if err != nil {
		return err
	}

	return nil
}

// SetScheduleStatus sets the status for the specified Phillips Hue schedule
func (h *Connection) SetScheduleStatus(schedule int, status string) error {
	// Error checking
	if !h.doesScheduleExist(schedule) {
		return fmt.Errorf("Schedule %d not found", schedule)
	}

	if strings.Trim(status, " ") != "enabled" && strings.Trim(status, " ") != "disabled" {
		return errors.New("Status must be one of the following: enabled, disabled")
	}

	attributes := fmt.Sprintf("{ \"status\": \"%s\" }", status)

	err := h.updateSchedule(schedule, attributes)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSchedule deletes the specified Phillips Hue schedule
func (h *Connection) DeleteSchedule(schedule int) error {
	// Error checking
	if !h.doesScheduleExist(schedule) {
		return fmt.Errorf("Schedule %d not found", schedule)
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/schedules/%d", h.baseURL, schedule), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dataStr := string(data)

	// Check for error in response
	if strings.Contains(dataStr, "error") {
		errMsg := dataStr[strings.Index(dataStr, "\"description\":\"")+15 : strings.Index(dataStr, "\"}}]")]
		return errors.New(errMsg)
	}

	return nil
}

func (h *Connection) doesScheduleExist(schedule int) bool {
	// If GetSchedule returns an error, then the schedule doesn't exist
	_, err := h.GetSchedule(schedule)
	if err != nil {
		return false
	}

	return true
}

func (h *Connection) updateSchedule(schedule int, attributes string) error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := strings.NewReader(attributes)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/schedules/%d", h.baseURL, schedule), body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dataStr := string(data)

	if strings.Contains(dataStr, "error") {
		errMsg := dataStr[strings.Index(dataStr, "\"description\":\"")+15 : strings.Index(dataStr, "\"}}]")]
		return errors.New(errMsg)
	}

	return nil
}

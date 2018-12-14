package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type ScheduleCommandBody struct {
	Scene string `json:"scene"`
	Flag  bool   `json:"flag"`
	On    bool   `json:"on"`
}

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

	fullResponse := string(body)

	schedule := Schedule{}
	allSchedules := []Schedule{}
	fullResponse = strings.Replace(fullResponse, "{", "", 1)

	// Get starting schedule id
	count := -1
	if len(fullResponse) > 0 {
		countStr := fullResponse[0:strings.Index(fullResponse, ":")]
		countStr = strings.Replace(countStr, "\"", "", -1)
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return nil, err
		}
	}

	// Format output
	for count != -1 {
		tmpArray := strings.Split(fullResponse, fmt.Sprintf("\"%d\":", count))

		if len(tmpArray) <= 1 {
			if len(tmpArray) > 0 {
				if tmpArray[0] != "" {
					//Remove leading or trailing commas
					tmpArray[0] = strings.Trim(tmpArray[0], ",")

					//If sting ends in two curly braces remove one
					if strings.LastIndex(tmpArray[0], "}}") == len(tmpArray[0])-2 {
						tmpArray[0] = tmpArray[0][0 : len(tmpArray[0])-1]
					}

					err = json.Unmarshal([]byte(tmpArray[0]), &schedule)
					if err != nil {
						return nil, err
					}

					schedule.ID = count - 1

					allSchedules = append(allSchedules, schedule)
				}
			}
			count = -1
		} else {
			if tmpArray[0] != "" {
				// Remove leading or trailing commas
				tmpArray[0] = strings.Trim(tmpArray[0], ",")

				err = json.Unmarshal([]byte(tmpArray[0]), &schedule)
				if err != nil {
					return nil, err
				}

				schedule.ID = count - 1

				allSchedules = append(allSchedules, schedule)
			}

			fullResponse = strings.Replace(fullResponse, fmt.Sprintf("\"%d\":", count), "", 1)
			fullResponse = strings.Replace(fullResponse, tmpArray[0], "", 1)
			count++
		}
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

// formatStruct formats a struct as a JSON string
func (h *Connection) formatStruct(data interface{}) string {
	str := "{"
	d := reflect.ValueOf(data)
	t := d.Type()

	for i := 0; i < d.NumField(); i++ {
		str += fmt.Sprintf("\"%s\": ", strings.ToLower(t.Field(i).Name))

		switch d.Field(i).Kind() {
		case reflect.String:
			str += fmt.Sprintf("\"%s\",", d.Field(i).Interface())
		// Check for slice and array
		case reflect.Struct:
			str += fmt.Sprintf("%s,", h.formatStruct(d.Field(i).Interface()))
		default:
			str += fmt.Sprintf("%v,", d.Field(i).Interface())
		}
	}
	str = str[:len(str)-1]
	str += "}"

	return str
}

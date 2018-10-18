package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type groupAction struct {
	On        bool      `json:"on"`
	Bri       int       `json:"bri"`
	Hue       int       `json:"hue"`
	Sat       int       `json:"sat"`
	Effect    string    `json:"effect"`
	XY        []float32 `json:"xy"`
	CT        int       `json:"ct"`
	Alert     string    `json:"alert"`
	ColorMode string    `json:"colormode"`
}

// Group contains all data returned from the Phillips Hue API
// for an individual Phillips Hue light group
type Group struct {
	Name   string      `json:"name"`
	Lights []string    `json:"lights"`
	Type   string      `json:"type"`
	Action groupAction `json:"action"`
	ID     int         `json:"id"`
}

// GetAllGroups gets all Phillips Hue light groups connected to current bridge
func (h *Connection) GetAllGroups() ([]Group, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/groups", h.baseURL))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fullResponse := string(body)

	group := Group{}
	allGroups := []Group{}
	fullResponse = strings.Replace(fullResponse, "{", "", 1)

	// Get starting light id
	count := -1
	if len(fullResponse) > 0 {
		countStr := fullResponse[0:strings.Index(fullResponse, ":")]
		countStr = strings.Replace(countStr, "\"", "", -1)
		count, err = strconv.Atoi(countStr)
		if err != nil {
			return nil, err
		}
	}

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

					err = json.Unmarshal([]byte(tmpArray[0]), &group)
					if err != nil {
						return nil, err
					}

					group.ID = count - 1

					allGroups = append(allGroups, group)
				}
			}
			count = -1
		} else {
			if tmpArray[0] != "" {
				// Remove leading or trailing commas
				tmpArray[0] = strings.Trim(tmpArray[0], ",")

				// If sting ends in two curly braces remove one
				if strings.LastIndex(tmpArray[0], "}}") == len(tmpArray[0])-2 {
					tmpArray[0] = tmpArray[0][0 : len(tmpArray[0])-1]
				}

				err = json.Unmarshal([]byte(tmpArray[0]), &group)
				if err != nil {
					return nil, err
				}

				group.ID = count - 1

				allGroups = append(allGroups, group)
			}

			fullResponse = strings.Replace(fullResponse, fmt.Sprintf("\"%d\":", count), "", 1)
			fullResponse = strings.Replace(fullResponse, tmpArray[0], "", 1)
			count++
		}
	}

	return allGroups, nil
}

// CreateGroup creates a new group with the specified name consisting of the specified
// lights. The group is added to the bridge using the next available ID.
func (h *Connection) CreateGroup(name, groupType string, lights []string) error {
	// POST - %s/groups
	// Body example -
	/*
		{
			"lights": ["1", "2"],
			"name": "bedroom",
			"type": "LightGroup"
		}
		OR if type is Room:
		{
			"lights": ["1", "2"],
			"name": "Living room",
			"type": "Room",
			"class": "Living room"
		}
	*/
	// Response example - [{"success": {"id": "1"}}]

	return nil
}

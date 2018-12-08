package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type groupState struct {
	AllOn bool `json:"all_on"`
	AnyOn bool `json:"any_on"`
}

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
	Name    string      `json:"name"`
	Lights  []string    `json:"lights"`
	Sensors []string    `json:"sensors"`
	Type    string      `json:"type"`
	State   groupState  `json:"state"`
	Recycle bool        `json:"recycle"`
	Action  groupAction `json:"action"`
	ID      int         `json:"id"`
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
func (h *Connection) CreateGroup(name, groupType, class string, lights []string) error {
	// Error checking
	name = strings.Trim(name, " ")
	if name == "" {
		return fmt.Errorf("Name must not be empty")
	}

	// LightGroup is the default group
	groupType = strings.Trim(groupType, " ")
	if groupType == "" {
		groupType = "LightGroup"
	}

	// LightGroup, Room, Luminaire, and LightSource = are valid groups
	if groupType != "LightGroup" && groupType != "Room" && groupType != "Luminaire" && groupType != "LightSource" {
		return fmt.Errorf("Group Type must be one of the following: LightGroup, Room, Luminaire, LightSource")
	}

	// Other is the default class
	class = strings.Trim(class, " ")
	if class == "" {
		class = "Other"
	}

	// Check that all lights are valid
	for idx, light := range lights {
		lightNum, err := strconv.Atoi(light)
		if err != nil {
			return fmt.Errorf("Light %s not found", light)
		}

		if !h.doesLightExist(lightNum) {
			return fmt.Errorf("Light %s not found", light)
		}

		// Add quotes around light number
		lights[idx] = fmt.Sprintf("\"%s\"", lights[idx])
	}

	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	reqBody := strings.NewReader(fmt.Sprintf("{\"name\": \"%s\", \"type\": \"%s\", \"class\": \"%s\", \"lights\": %v}", name, groupType, class, lights))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fullResponse := string(body)
	fullResponse = strings.ToLower(fullResponse)

	if !strings.Contains(fullResponse, "success") {
		errMsg := fullResponse[strings.Index(fullResponse, "description")+14 : strings.LastIndex(fullResponse, "\"")]
		return fmt.Errorf("Unable to create group %s: %s", name, errMsg)
	}

	return nil
}

// GetGroup gets the specified Phillips Hue light group
func (h *Connection) GetGroup(group int) (Group, error) {
	err := h.initializeHue()
	if err != nil {
		return Group{}, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/groups/%d", h.baseURL, group))
	if err != nil {
		return Group{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Group{}, err
	}

	// Group not found
	if len(body) == 0 {
		return Group{}, fmt.Errorf("Group not found")
	}

	groupRes := Group{}

	err = json.Unmarshal(body, &groupRes)
	if err != nil {
		return Group{}, err
	}

	return groupRes, nil
}

func (h *Connection) DeleteGroup(group int) error {
	// Error checking:
	// - Does group exist
	// - Group type must not be LightSource or Luminaire

	// DELETE - /api/<username>/groups/<id>

	// Sample Response:
	/*
		[{
			"success": "/groups/1 deleted."
		}]
	*/

	return nil
}

func (h *Connection) doesGroupExist(group int) bool {

	return true
}

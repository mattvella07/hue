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
func (h *Connection) CreateGroup(name, groupType, class string, lights []int) error {
	// Error checking
	name = strings.Trim(name, " ")
	if name == "" {
		return errors.New("Name must not be empty")
	}

	// LightGroup is the default group
	groupType = strings.Trim(groupType, " ")
	if groupType == "" {
		groupType = "LightGroup"
	}

	// LightGroup, Room, Luminaire, and LightSource are valid groups
	if groupType != "LightGroup" && groupType != "Room" && groupType != "Luminaire" && groupType != "LightSource" {
		return errors.New("Group Type must be one of the following: LightGroup, Room, Luminaire, LightSource")
	}

	// Other is the default class
	class = strings.Trim(class, " ")
	if class == "" {
		class = "Other"
	}

	if !h.allLightsValid(lights) {
		return errors.New("One of the lights is invalid")
	}

	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	reqBody := strings.NewReader(fmt.Sprintf("{\"name\": \"%s\", \"type\": \"%s\", \"class\": \"%s\", \"lights\": %s}", name, groupType, class, h.formatSlice(lights)))
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
		return Group{}, errors.New("Group not found")
	}

	groupRes := Group{}

	err = json.Unmarshal(body, &groupRes)
	if err != nil {
		return Group{}, err
	}

	return groupRes, nil
}

// RenameGroup renames the specified Phillips Hue group
func (h *Connection) RenameGroup(group int, name string) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	attributes := fmt.Sprintf("{ \"name\": \"%s\" }", name)

	err := h.changeGroupAttributes(group, attributes)
	if err != nil {
		return err
	}

	return nil
}

// SetLightsInGroup sets the lights that are in the specified Phillips Hue group
func (h *Connection) SetLightsInGroup(group int, lights []int) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	if !h.allLightsValid(lights) {
		return errors.New("One of the lights is invalid")
	}

	attributes := fmt.Sprintf("{ \"lights\": %s }", h.formatSlice(lights))

	err := h.changeGroupAttributes(group, attributes)
	if err != nil {
		return err
	}

	return nil
}

// SetGroupClass sets the class for the specified Phillips Hue group
func (h *Connection) SetGroupClass(group int, class string) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	if strings.Trim(class, " ") == "" {
		return errors.New("Class must not be empty")
	}

	attributes := fmt.Sprintf("{ \"class\": \"%s\" }", class)

	err := h.changeGroupAttributes(group, attributes)
	if err != nil {
		return err
	}

	return nil
}

// 2.5 Set Group State

// DeleteGroup deletes the specified Phillips Hue light group
func (h *Connection) DeleteGroup(group int) error {
	// Error checking
	currentGroup, err := h.GetGroup(group)
	if err != nil {
		return fmt.Errorf("Group %d not found", group)
	}

	// Groups with type LightSource or Luminaire can't be deleted
	if currentGroup.Type == "LightSource" || currentGroup.Type == "Luminaire" {
		return fmt.Errorf("Unable to delete group %d: Can't delete group with a type of LightSource or Luminaire", group)
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%d", h.baseURL, group), nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) doesGroupExist(group int) bool {
	// If GetGroup returns an error, then the group doesn't exist
	_, err := h.GetGroup(group)
	if err != nil {
		return false
	}

	return true
}

func (h *Connection) allLightsValid(lights []int) bool {
	for _, light := range lights {
		if !h.doesLightExist(light) {
			return false
		}
	}

	return true
}

func (h *Connection) formatSlice(sli []int) string {
	str := "["
	for _, s := range sli {
		str += fmt.Sprintf("\"%d\",", s)
	}
	str = str[:len(str)-1]
	str += "]"

	return str
}

func (h *Connection) changeGroupAttributes(group int, attributes string) error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := strings.NewReader(attributes)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%d", h.baseURL, group), body)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
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

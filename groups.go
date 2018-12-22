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
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return []Group{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(body, &fullResponse)
	if err != nil {
		return nil, err
	}

	allGroups := []Group{}

	// Loop through all keys in the map and Unmarshal into
	// Group type
	for key, val := range fullResponse {
		group := Group{}

		g, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(g, &group)
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}

		group.ID = id

		allGroups = append(allGroups, group)
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
	defer resp.Body.Close()

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

	err := h.updateGroup(group, "attributes", attributes)
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

	err := h.updateGroup(group, "attributes", attributes)
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

	err := h.updateGroup(group, "attributes", attributes)
	if err != nil {
		return err
	}

	return nil
}

// TurnOnAllLightsInGroup turns on all lights in the specified Phillips Hue group
// without setting the color
func (h *Connection) TurnOnAllLightsInGroup(group int) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	state := "{ \"on\": true }"

	err := h.updateGroup(group, "state", state)
	if err != nil {
		return err
	}

	return nil
}

// TurnOnAllLightsInGroupWithColor turns on all lights in the specified Phillips Hue group
// to the color specified by the x and y parameters. Also sets the Bri, Hue, and Sat
// properties
func (h *Connection) TurnOnAllLightsInGroupWithColor(group int, x, y float32, bri, hue, sat int) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	err := h.validateColorParams(x, y, bri, hue, sat)
	if err != nil {
		return err
	}

	state := fmt.Sprintf("{\"on\": true, \"xy\": [%f, %f], \"bri\": %d, \"hue\": %d, \"sat\": %d}", x, y, bri, hue, sat)

	err = h.updateGroup(group, "state", state)
	if err != nil {
		return err
	}

	return nil
}

// TurnOffAllLightsInGroup turns off all lights in the specified Phillips Hue group
func (h *Connection) TurnOffAllLightsInGroup(group int) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	state := "{ \"on\": false }"

	err := h.updateGroup(group, "state", state)
	if err != nil {
		return err
	}

	return nil
}

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

func (h *Connection) doesGroupExist(group int) bool {
	// If GetGroup returns an error, then the group doesn't exist
	_, err := h.GetGroup(group)
	if err != nil {
		return false
	}

	return true
}

func (h *Connection) updateGroup(group int, toUpdate, value string) error {
	url := ""
	switch toUpdate {
	case "attributes":
		url = fmt.Sprintf("%s/groups/%d", h.baseURL, group)
	case "state":
		url = fmt.Sprintf("%s/groups/%d/action", h.baseURL, group)
	default:
		return fmt.Errorf("Error while updating group %d", group)
	}

	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := strings.NewReader(value)
	req, err := http.NewRequest("PUT", url, body)
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

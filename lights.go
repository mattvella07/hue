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

type lightState struct {
	On        bool      `json:"on"`
	Bri       int       `json:"bri"`
	Hue       int       `json:"hue"`
	Sat       int       `json:"sat"`
	Effect    string    `json:"effect"`
	XY        []float32 `json:"xy"`
	CT        int       `json:"ct"`
	Alert     string    `json:"alert"`
	ColorMode string    `json:"colormode"`
	Mode      string    `json:"mode"`
	Reachable bool      `json:"reachable"`
}

type lightSWUpdate struct {
	State       string `json:"state"`
	LastInstall string `json:"lastinstall"`
}

type lightCapabilitiesCT struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type lightCapabilitiesControl struct {
	MindimLevel    int                 `json:"mindimlevel"`
	MaxLumen       int                 `json:"maxlumen"`
	ColorGamutType string              `json:"colorgamuttype"`
	ColorGamut     [][]float32         `json:"colorgamut"`
	CT             lightCapabilitiesCT `json:"ct"`
}

type lightCapabilitiesStreaming struct {
	Renderer bool `json:"renderer"`
	Proxy    bool `json:"proxy"`
}

type lightCapabilities struct {
	Certified bool                       `json:"certified"`
	Control   lightCapabilitiesControl   `json:"control"`
	Streaming lightCapabilitiesStreaming `json:"streaming"`
}

type lightConfig struct {
	ArcheType string `json:"archetype"`
	Function  string `json:"function"`
	Direction string `json:"direction"`
}

// Light contains all data returned from the Phillips Hue API
// for an individual Phillips Hue light
type Light struct {
	State            lightState        `json:"state"`
	SWUpdate         lightSWUpdate     `json:"swupdate"`
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	ModelID          string            `json:"modelid"`
	ManufacturerName string            `json:"manufacturername"`
	ProductName      string            `json:"productname"`
	Capabilities     lightCapabilities `json:"capabilities"`
	Config           lightConfig       `json:"config"`
	UniqueID         string            `json:"uniqueid"`
	SWVersion        string            `json:"swversion"`
	SWConfigID       string            `json:"swconfigid"`
	ProductID        string            `json:"productid"`
	ID               int               `json:"id"`
}

// NewLight contains all data for a new Phillips Hue light
type NewLight struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// NewLightResponse is the response from the Phillips Hue API for
// all new lights
type NewLightResponse struct {
	NewLights []NewLight `json:"newLights"`
	LastScan  string     `json:"lastScan"`
}

// GetLights gets all Phillips Hue lights connected to current bridge
func (h *Connection) GetLights() ([]Light, error) {
	body, err := h.get("lights")
	if err != nil {
		return []Light{}, err
	}

	if len(body) == 0 {
		return []Light{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(body, &fullResponse)
	if err != nil {
		return nil, err
	}

	allLights := []Light{}

	// Loop through all keys in the map and Unmarshal into
	// Group type
	for key, val := range fullResponse {
		light := Light{}

		l, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(l, &light)
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}

		light.ID = id

		allLights = append(allLights, light)
	}

	return allLights, nil
}

// GetNewLights gets Phillips Hue lights that were discovered the last time
// FindNewLights was called
func (h *Connection) GetNewLights() (NewLightResponse, error) {
	body, err := h.get("lights/new")
	if err != nil {
		return NewLightResponse{}, err
	}

	if len(body) == 0 {
		return NewLightResponse{}, nil
	}

	fullResponse := string(body)

	if len(fullResponse) <= 1 {
		return NewLightResponse{}, errors.New("No new lights found")
	}

	// Remove first and last character - { and }
	fullResponse = fullResponse[1 : len(fullResponse)-1]

	// Split based on ,
	fullResponseArr := strings.Split(fullResponse, ",")

	newLight := NewLight{}
	newLightRes := NewLightResponse{}

	// Loop through response and transform into type NewLightResponse
	for _, item := range fullResponseArr {
		if strings.Contains(item, "\"lastscan\"") {
			// Remove text lastscan and " characters
			newLightRes.LastScan = strings.Replace(item, "\"lastscan\":", "", 1)
			newLightRes.LastScan = strings.Replace(newLightRes.LastScan, "\"", "", -1)
		} else {
			// Must be at least 3 items after splitting by :
			items := strings.Split(item, ":")
			if len(items) < 3 {
				return NewLightResponse{}, errors.New("Error processing result")
			}

			// Remove " characters to get ID
			if id, err := strconv.Atoi(strings.Replace(items[0], "\"", "", -1)); err == nil {
				newLight.ID = id
			} else {
				return NewLightResponse{}, errors.New("Error processing result")
			}

			// Remove " and } characters to get light name
			newLight.Name = strings.Replace(items[2], "\"", "", -1)
			newLight.Name = strings.Replace(newLight.Name, "}", "", -1)

			newLightRes.NewLights = append(newLightRes.NewLights, newLight)
		}
	}

	return newLightRes, nil
}

// FindNewLights finds new Phillips Hue lights that have been added since
// the last time performing this call
func (h *Connection) FindNewLights() error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/lights", h.baseURL), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetLight gets the specified Phillips Hue light
func (h *Connection) GetLight(light int) (Light, error) {
	body, err := h.get(fmt.Sprintf("lights/%d", light))
	if err != nil {
		return Light{}, err
	}

	// Light not found
	if len(body) == 0 {
		return Light{}, fmt.Errorf("Light %d not found", light)
	}

	lightRes := Light{}

	err = json.Unmarshal(body, &lightRes)
	if err != nil {
		return Light{}, err
	}

	return lightRes, nil
}

// RenameLight renames the specified Phillips Hue light
func (h *Connection) RenameLight(light int, name string) error {
	// Error checking
	if !h.doesLightExist(light) {
		return fmt.Errorf("Light %d not found", light)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	client := &http.Client{}
	reqBody := strings.NewReader(fmt.Sprintf("{ \"name\": \"%s\" }", name))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/lights/%d", h.baseURL, light), reqBody)
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

	if !strings.Contains(fullResponse, "updated") && !strings.Contains(fullResponse, "success") {
		return fmt.Errorf("Unable to rename light %d to %s", light, name)
	}

	return nil
}

// TurnOnLight turns on the specified Phillips Hue light without setting the color
func (h *Connection) TurnOnLight(light int) error {
	// Error checking
	if !h.doesLightExist(light) {
		return fmt.Errorf("Light %d not found", light)
	}

	// Set state
	state := "{\"on\": true}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

// TurnOnLightWithColor turns on the specified Phillips Hue light to the color
// specified by the x and y parameters. Also sets the Bri, Hue, and Sat properties
func (h *Connection) TurnOnLightWithColor(light int, x, y float32, bri, hue, sat int) error {
	// Error checking
	if !h.doesLightExist(light) {
		return fmt.Errorf("Light %d not found", light)
	}

	err := h.validateColorParams(x, y, bri, hue, sat)
	if err != nil {
		return err
	}

	// Set state
	state := fmt.Sprintf("{\"on\": true, \"xy\": [%f, %f], \"bri\": %d, \"hue\": %d, \"sat\": %d}", x, y, bri, hue, sat)

	err = h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

// TurnOffLight turns off the specified Phillips Hue light
func (h *Connection) TurnOffLight(light int) error {
	// Error checking
	if !h.doesLightExist(light) {
		return fmt.Errorf("Light %d not found", light)
	}

	// Set state
	state := "{\"on\": false}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

// DeleteLight deletes a Phillips Hue light from the bridge
func (h *Connection) DeleteLight(light int) error {
	// Error checking
	if !h.doesLightExist(light) {
		return fmt.Errorf("Light %d not found", light)
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/lights/%d", h.baseURL, light), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) doesLightExist(light int) bool {
	// If GetLight returns an error, then the light doesn't exist
	_, err := h.GetLight(light)
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

func (h *Connection) changeLightState(light int, state string) error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := strings.NewReader(state)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/lights/%d/state", h.baseURL, light), body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (h *Connection) validateColorParams(x, y float32, bri, hue, sat int) error {
	if x < 0 || x > 1 {
		return errors.New("Invalid color value: x must be between 0 and 1")
	}

	if y < 0 || y > 1 {
		return errors.New("Invalid color value: y must be between 0 and 1")
	}

	if bri < 1 || bri > 254 {
		return errors.New("Invalid brightness value: bri must be between 1 and 254")
	}

	if hue < 0 || hue > 65535 {
		return errors.New("Invalid hue value: hue must be between 0 and 65,535")
	}

	if sat < 0 || sat > 254 {
		return errors.New("Invalid saturation value: sat must be between 0 and 254")
	}

	return nil
}

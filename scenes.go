package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SceneAppData contains the data for the AppData field in the
// Scene type
type SceneAppData struct {
	Version int    `json:"version"`
	Data    string `json:"data"`
}

// Scene contains all data returned from the Phillips Hue API
// for an individual Phillips Hue scene
type Scene struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Group       string       `json:"group"`
	Lights      []string     `json:"lights"`
	Owner       string       `json:"owner"`
	Recycle     bool         `json:"recycle"`
	Locked      bool         `json:"locked"`
	AppData     SceneAppData `json:"appdata"`
	Picture     string       `json:"picture"`
	LastUpdated string       `json:"lastupdated"`
	Version     int          `json:"version"`
	ID          string       `json:"id"`
}

// GetScenes gets all Phillips Hue scenes
func (h *Connection) GetScenes() ([]Scene, error) {
	data, err := h.get("scenes")
	if err != nil {
		return []Scene{}, err
	}

	if len(data) == 0 {
		return []Scene{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return []Scene{}, err
	}

	allScenes := []Scene{}

	// Loop through all keys in the map and Unmarshal into
	// Scene type
	for key, val := range fullResponse {
		scene := Scene{}

		s, err := json.Marshal(val)
		if err != nil {
			return []Scene{}, err
		}

		err = json.Unmarshal(s, &scene)
		if err != nil {
			return []Scene{}, err
		}

		scene.ID = key

		allScenes = append(allScenes, scene)
	}

	return allScenes, nil
}

// CreateLightScene creates a new scene of type LightScene with the specified name
func (h *Connection) CreateLightScene(name string, lights []int, recycle bool, appData SceneAppData) error {
	// Error checking
	if len(lights) == 0 {
		return errors.New("Lights must not be empty")
	}

	if !h.allLightsValid(lights) {
		return errors.New("One of the lights is invalid")
	}

	bodyStr := fmt.Sprintf("{\"name\": \"%s\", \"type\": \"LightScene\", \"lights\": %s, \"recycle\": %t", name, h.formatSlice(lights), recycle)

	if appData.Version != 0 || strings.Trim(appData.Data, " ") != "" {
		bodyStr += fmt.Sprintf(", \"appdata\": %s", h.formatStruct(appData))
	}
	bodyStr += "}"

	reqBody := strings.NewReader(bodyStr)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/scenes", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateGroupScene creates a new scene of type GroupScene with the specified name
func (h *Connection) CreateGroupScene(name string, group int, recycle bool, appData SceneAppData) error {
	// Error checking
	if !h.doesGroupExist(group) {
		return fmt.Errorf("Group %d not found", group)
	}

	bodyStr := fmt.Sprintf("{\"name\": \"%s\", \"type\": \"GroupScene\", \"group\": \"%d\", \"recycle\": %t", name, group, recycle)

	if appData.Version != 0 || strings.Trim(appData.Data, " ") != "" {
		bodyStr += fmt.Sprintf(", \"appdata\": %s", h.formatStruct(appData))
	}
	bodyStr += "}"

	reqBody := strings.NewReader(bodyStr)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/scenes", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// RenameScene renames the specified Phillips Hue scene
func (h *Connection) RenameScene(scene, name string) error {
	// Error checking
	if !h.doesSceneExist(scene) {
		return fmt.Errorf("Scene %s not found", scene)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	attributes := fmt.Sprintf("{ \"name\": \"%s\" }", name)

	err := h.updateScene(scene, attributes)
	if err != nil {
		return err
	}

	return nil
}

// SetLightsInScene sets the lights that are in the specified Phillips Hue scene
func (h *Connection) SetLightsInScene(scene string, lights []int) error {
	// Error checking
	if !h.doesSceneExist(scene) {
		return fmt.Errorf("Scene %s not found", scene)
	}

	if len(lights) == 0 {
		return errors.New("Lights must not be empty")
	}

	if !h.allLightsValid(lights) {
		return errors.New("One of the lights is invalid")
	}

	attributes := fmt.Sprintf("{ \"lights\": %s }", h.formatSlice(lights))

	err := h.updateScene(scene, attributes)
	if err != nil {
		return err
	}

	return nil
}

// DeleteScene deletes the specified Phillips Hue scene
func (h *Connection) DeleteScene(scene string) error {
	// Error checking
	if !h.doesSceneExist(scene) {
		return fmt.Errorf("Scene %s not found", scene)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/scenes/%s", h.baseURL, scene), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// GetScene gets the specified Phillips Hue scene by ID
func (h *Connection) GetScene(scene string) (Scene, error) {
	data, err := h.get(fmt.Sprintf("scenes/%s", scene))
	if err != nil {
		return Scene{}, err
	}

	// Scene not found
	if len(data) == 0 {
		return Scene{}, fmt.Errorf("Scene %s not found", scene)
	}

	sceneRes := Scene{}

	err = json.Unmarshal(data, &sceneRes)
	if err != nil {
		return Scene{}, err
	}

	return sceneRes, nil
}

func (h *Connection) doesSceneExist(scene string) bool {
	// Scene ID must not be empty
	if strings.Trim(scene, " ") == "" {
		return false
	}

	// If GetScene returns an error, then the scene doesn't exist
	_, err := h.GetScene(scene)
	if err != nil {
		return false
	}

	return true
}

func (h *Connection) updateScene(scene, value string) error {
	url := fmt.Sprintf("%s/scenes/%s", h.baseURL, scene)

	reqBody := strings.NewReader(value)
	req, err := http.NewRequest("PUT", url, reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

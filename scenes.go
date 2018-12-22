package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type sceneAppData struct {
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
	AppData     sceneAppData `json:"appdata"`
	Picture     string       `json:"picture"`
	LastUpdated string       `json:"lastupdated"`
	Version     int          `json:"version"`
	ID          string       `json:"id"`
}

// GetAllScenes gets all Phillips Hue scenes
func (h *Connection) GetAllScenes() ([]Scene, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/scenes", h.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return []Scene{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(body, &fullResponse)
	if err != nil {
		return nil, err
	}

	allScenes := []Scene{}

	// Loop through all keys in the map and Unmarshal into
	// Scene type
	for key, val := range fullResponse {
		scene := Scene{}

		s, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(s, &scene)
		if err != nil {
			return nil, err
		}

		scene.ID = key

		allScenes = append(allScenes, scene)
	}

	return allScenes, nil
}

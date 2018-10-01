package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (h *Connection) GetAllLights() ([]hueLight, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/lights", h.baseURL))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fullResponse := string(body)

	light := hueLight{}
	count := 1
	fullResponse = strings.Replace(fullResponse, "{", "", 1)
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

					err = json.Unmarshal([]byte(tmpArray[0]), &light)
					if err != nil {
						return nil, err
					}

					h.Lights = append(h.Lights, light)
				}
			}
			count = -1
		} else {
			if tmpArray[0] != "" {
				//Remove leading or trailing commas
				tmpArray[0] = strings.Trim(tmpArray[0], ",")

				//If sting ends in two curly braces remove one
				if strings.LastIndex(tmpArray[0], "}}") == len(tmpArray[0])-2 {
					tmpArray[0] = tmpArray[0][0 : len(tmpArray[0])-1]
				}

				err = json.Unmarshal([]byte(tmpArray[0]), &light)
				if err != nil {
					return nil, err
				}

				h.Lights = append(h.Lights, light)
			}

			fullResponse = strings.Replace(fullResponse, fmt.Sprintf("\"%d\":", count), "", 1)
			fullResponse = strings.Replace(fullResponse, tmpArray[0], "", 1)
			count++
		}
	}

	return h.Lights, nil
}

func (h *Connection) GetLight(light int) (hueLight, error) {
	err := h.initializeHue()
	if err != nil {
		return hueLight{}, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/lights/%d", h.baseURL, light))
	if err != nil {
		return hueLight{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hueLight{}, err
	}

	lightRes := hueLight{}

	err = json.Unmarshal(body, &lightRes)
	if err != nil {
		return hueLight{}, err
	}

	return lightRes, nil
}

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

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
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

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (h *Connection) TurnOnLight(light int) error {
	state := "{\"on\": true}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) TurnOnLightWithColor(light int, x, y float32, bri, hue, sat int) error {
	state := fmt.Sprintf("{\"on\": true, \"xy\": [%f, %f], \"bri\": %d, \"hue\": %d, \"sat\": %d}", x, y, bri, hue, sat)

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) TurnOffLight(light int) error {
	state := "{\"on\": false}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

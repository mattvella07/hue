package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (h *Connection) GetLights() ([]hueLight, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(h.baseURL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fullResponse := string(body)

	var light hueLight
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

func (h *Connection) changeLightState(light int, state string) error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := strings.NewReader(state)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%d/state", h.baseURL, light), body)
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

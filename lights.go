package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetAllLights gets all Phillips Hue lights connected to current bridge
func (h *Connection) GetAllLights() ([]Light, error) {
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

	light := Light{}
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

// GetLight gets the specified Phillips Hue light
func (h *Connection) GetLight(light int) (Light, error) {
	err := h.initializeHue()
	if err != nil {
		return Light{}, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/lights/%d", h.baseURL, light))
	if err != nil {
		return Light{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Light{}, err
	}

	//Light not found
	if len(body) == 0 {
		return Light{}, fmt.Errorf("Light not found")
	}

	lightRes := Light{}

	err = json.Unmarshal(body, &lightRes)
	if err != nil {
		return Light{}, err
	}

	return lightRes, nil
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

// TurnOnLight turns on the specified Phillips Hue light without setting the color
func (h *Connection) TurnOnLight(light int) error {
	//Error checking - check light to make sure it exists in the Lights array

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
	//Error checking - check light to make sure it exists in the Lights array
	//bri - Between 1 and 254
	//hue - Between 0 and 65535
	//sat - Between 0 and 254

	state := fmt.Sprintf("{\"on\": true, \"xy\": [%f, %f], \"bri\": %d, \"hue\": %d, \"sat\": %d}", x, y, bri, hue, sat)

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

// TurnOffLight turns off the specified Phillips Hue light
func (h *Connection) TurnOffLight(light int) error {
	//Error checking - check light to make sure it exists in the Lights array

	state := "{\"on\": false}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

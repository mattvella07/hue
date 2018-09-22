package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Hue struct {
	discoveryResponse []hueDiscoveryResponse
	internalIPAddress string
	userID            string
	baseURL           string
	Lights            []hueLight
}

type hueDiscoveryResponse struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

type hueLight struct {
	State struct {
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
	} `json:"state"`
	SWUpdate struct {
		State       string `json:"state"`
		LastInstall string `json:"lastinstall"`
	} `json:"swupdate"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	ModelID          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	ProductName      string `json:"productname"`
	Capabilities     struct {
		Certified bool `json:"certified"`
		Control   struct {
			MindimLevel    int         `json:"mindimlevel"`
			MaxLumen       int         `json:"maxlumen"`
			ColorGamutType string      `json:"colorgamuttype"`
			ColorGamut     [][]float32 `json:"colorgamut"`
			CT             struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"ct"`
		} `json:"control"`
		Streaming struct {
			Renderer bool `json:"renderer"`
			Proxy    bool `json:"proxy"`
		} `json:"streaming"`
	} `json:"capabilities"`
	Config struct {
		ArcheType string `json:"archetype"`
		Function  string `json:"function"`
		Direction string `json:"direction"`
	} `json:"config"`
	UniqueID   string `json:"uniqueid"`
	SWVersion  string `json:"swversion"`
	SWConfigID string `json:"swconfigid"`
	ProductID  string `json:"productid"`
}

const hueDiscoveryURL = "https://discovery.meethue.com/"

func (h *Hue) initializeHue() error {
	var err error

	if h.internalIPAddress == "" {
		err = h.getBridgeIPAddress()
		if err != nil {
			return fmt.Errorf("GetBridgeIPAddress Error: %s", err)
		}
	}

	if h.userID == "" {
		err = h.getUserID()
		if err != nil {
			return fmt.Errorf("GetUserID Error: %s", err)
		}
	}

	if h.baseURL == "" {
		h.getBaseURL()
	}

	return nil
}

func (h *Hue) getBridgeIPAddress() error {
	resp, err := http.Get(hueDiscoveryURL)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &h.discoveryResponse)
	if err != nil {
		return err
	}

	if len(h.discoveryResponse) == 0 {
		return errors.New("Unable to determine Hue bridge internal IP address")
	}

	h.internalIPAddress = h.discoveryResponse[0].InternalIPAddress
	return nil
}

func (h *Hue) getUserID() error {
	val, ok := os.LookupEnv("hueUserID")
	if ok {
		h.userID = val
		return nil
	} else {
		return errors.New("Unable to get Hue user ID")
		//Generate it and set env var
		//fmt.Sprintf("http://%s/api", h.internalIPAddress)
	}
}

func (h *Hue) getBaseURL() {
	h.baseURL = fmt.Sprintf("http://%s/api/%s/lights/", h.internalIPAddress, h.userID)
}

func (h *Hue) GetLights() ([]hueLight, error) {
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

func (h *Hue) changeLightState(light int, state string) error {
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

func (h *Hue) TurnOnLight(light int) error {
	state := "{\"on\": true}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hue) TurnOffLight(light int) error {
	state := "{\"on\": false}"

	err := h.changeLightState(light, state)
	if err != nil {
		return err
	}

	return nil
}

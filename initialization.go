package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Connection struct {
	discoveryResponse []hueDiscoveryResponse
	internalIPAddress string
	UserID            string
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

func (h *Connection) initializeHue() error {
	var err error

	if h.internalIPAddress == "" {
		err = h.getBridgeIPAddress()
		if err != nil {
			return fmt.Errorf("GetBridgeIPAddress Error: %s", err)
		}
	}

	if h.UserID == "" {
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

func (h *Connection) getBridgeIPAddress() error {
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

func (h *Connection) getUserID() error {
	val, ok := os.LookupEnv("hueUserID")
	if ok {
		h.UserID = val
		return nil
	} else {
		return errors.New("Unable to get Hue user ID")
		//Generate it and set env var
		//fmt.Sprintf("http://%s/api", h.internalIPAddress)
	}
}

func (h *Connection) getBaseURL() {
	h.baseURL = fmt.Sprintf("http://%s/api/%s/lights/", h.internalIPAddress, h.UserID)
}

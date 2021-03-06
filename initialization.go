package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Connection contains important connection info
type Connection struct {
	discoveryResponse []hueDiscoveryResponse
	internalIPAddress string
	UserID            string
	baseURL           string
	isInitialized     bool
}

type hueDiscoveryResponse struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

const hueDiscoveryURL = "https://discovery.meethue.com/"

func (h *Connection) initializeHue() error {
	if h.isInitialized {
		return nil
	}

	err := h.getBridgeIPAddress()
	if err != nil {
		return fmt.Errorf("GetBridgeIPAddress Error: %s", err)
	}

	h.getBaseURL()

	h.isInitialized = true

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

func (h *Connection) getBaseURL() {
	h.baseURL = fmt.Sprintf("http://%s/api/%s", h.internalIPAddress, h.UserID)
}

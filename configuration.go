package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ConfigurationPortalState contains the data for the PortalState field in the
// Configuration type
type ConfigurationPortalState struct {
	SignedOn      bool   `json:"signedon"`
	Incoming      bool   `json:"incoming"`
	Outgoing      bool   `json:"outgoing"`
	Communication string `json:"communication"`
}

// ConfigurationSWUpdate contains the data for the SWUpdate field in the
// Configuration type
type ConfigurationSWUpdate struct {
	UpdateState int    `json:"updatestate"`
	URL         string `json:"url"`
	Text        string `json:"text"`
	Notify      bool   `json:"notify"`
}

// ConfigurationWhitelist contains the data for the Whitelist field in the
// Configuration type
type ConfigurationWhitelist struct {
	LastUseDate string `json:"last use date"`
	CreateDate  string `json:"create date"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

// Configuration contains all data returned from the Phillips Hue API
// for the Phillips Hue configuration
type Configuration struct {
	Name             string                   `json:"name"`
	ZigbeeChannel    int                      `json:"zigbeechannel"`
	Mac              string                   `json:"mac"`
	DHCP             bool                     `json:"dhcp"`
	IPAddress        string                   `json:"ipaddress"`
	NetMask          string                   `json:"netmask"`
	Gateway          string                   `json:"gateway"`
	ProxyAddress     string                   `json:"proxyaddress"`
	ProxyPort        int                      `json:"proxyport"`
	UTC              string                   `json:"UTC"`
	LocalTime        string                   `json:"localtime"`
	Timezone         string                   `json:"timezone"`
	Whitelist        []ConfigurationWhitelist `json:"whitelist"`
	SWVersion        string                   `json:"swversion"`
	APIVersion       string                   `json:"apiversion"`
	SWUpdate         ConfigurationSWUpdate    `json:"swupdate"`
	LinkButton       bool                     `json:"linkbutton"`
	PortalServices   bool                     `json:"portalservices"`
	PortalConnection string                   `json:"portalconnection"`
	PortalState      ConfigurationPortalState `json:"portalstate"`
}

// CreateUser creates a new user
func (h *Connection) CreateUser(deviceType string) error {
	// Error checking
	if strings.Trim(deviceType, " ") == "" {
		return errors.New("deviceType must not be empty")
	}

	reqBody := strings.NewReader(fmt.Sprintf("{\"devicetype\": \"%s\"}", deviceType))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// GetConfiguration gets the Phillips Hue configuration
func (h *Connection) GetConfiguration() (Configuration, error) {
	data, err := h.get("config")
	if err != nil {
		return Configuration{}, err
	}

	if len(data) == 0 {
		return Configuration{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return Configuration{}, err
	}

	config := Configuration{}

	// Create map to store Whitelist JSON response
	whitelist := make(map[string]interface{})

	whitelistData, err := json.Marshal(fullResponse["whitelist"])
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(whitelistData, &whitelist)
	if err != nil {
		return Configuration{}, err
	}

	// Loop through Whitelist slice and add to config
	for key, val := range whitelist {
		configWhitelist := ConfigurationWhitelist{}

		w, err := json.Marshal(val)
		if err != nil {
			return Configuration{}, err
		}

		err = json.Unmarshal(w, &configWhitelist)
		if err != nil {
			return Configuration{}, err
		}

		configWhitelist.ID = key

		config.Whitelist = append(config.Whitelist, configWhitelist)
	}

	// Remove whitelist data from JSON response so the rest can be
	// Unmarshalled into a Configuration object
	delete(fullResponse, "whitelist")

	otherData, err := json.Marshal(fullResponse)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(otherData, &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}

// DeleteUser deletes the specified user from the whitelist
func (h *Connection) DeleteUser(user string) error {
	// Error checking
	if strings.Trim(user, " ") == "" {
		return errors.New("User must not be empty")
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/config/whitelist/%s", h.baseURL, user), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

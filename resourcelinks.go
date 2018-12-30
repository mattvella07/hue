package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ResourceLink contains all data returned from the Phillips Hue API
// for an individual Phillips Hue resource link
type ResourceLink struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Class       int      `json:"class"`
	Owner       string   `json:"owner"`
	Links       []string `json:"links"`
	ID          int      `json:"id"`
}

// GetResourceLinks gets all Phillips Hue resource links
func (h *Connection) GetResourceLinks() ([]ResourceLink, error) {
	data, err := h.get("resourcelinks")
	if err != nil {
		return []ResourceLink{}, err
	}

	if len(data) == 0 {
		return []ResourceLink{}, nil
	}

	// Create map to store JSON response
	fullResponse := make(map[string]interface{})

	err = json.Unmarshal(data, &fullResponse)
	if err != nil {
		return []ResourceLink{}, err
	}

	allResourcelinks := []ResourceLink{}

	// Loop through all keys in the map and Unmarshal into
	// ResourceLink type
	for key, val := range fullResponse {
		link := ResourceLink{}

		s, err := json.Marshal(val)
		if err != nil {
			return []ResourceLink{}, err
		}

		err = json.Unmarshal(s, &link)
		if err != nil {
			return []ResourceLink{}, err
		}

		id, err := strconv.Atoi(key)
		if err != nil {
			return []ResourceLink{}, err
		}

		link.ID = id

		allResourcelinks = append(allResourcelinks, link)
	}

	return allResourcelinks, nil
}

// GetResourceLink gets the specified Phillips Hue resource link
func (h *Connection) GetResourceLink(resourceLink int) (ResourceLink, error) {
	data, err := h.get(fmt.Sprintf("resourcelinks/%d", resourceLink))
	if err != nil {
		return ResourceLink{}, err
	}

	// Resource link not found
	if len(data) == 0 {
		return ResourceLink{}, fmt.Errorf("Resource link %d not found", resourceLink)
	}

	linkRes := ResourceLink{}

	err = json.Unmarshal(data, &linkRes)
	if err != nil {
		return ResourceLink{}, err
	}

	return linkRes, nil
}

// CreateResourceLink creates a new resource link with the specified name
func (h *Connection) CreateResourceLink(name, description string, recycle bool, links []string) error {
	// Error checking
	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	if len(links) == 0 {
		return errors.New("Links must not be empty")
	}

	reqBody := strings.NewReader(fmt.Sprintf("{\"name\": \"%s\", \"description\": \"%s\", \"recycle\": %t, \"links\": %s}", name, description, recycle, links))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/resourcelinks", h.baseURL), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// RenameResourceLink renames the specified Phillips Hue resource link
func (h *Connection) RenameResourceLink(resourceLink int, name string) error {
	// Error checking
	if !h.doesResourceLinkExist(resourceLink) {
		return fmt.Errorf("Resource link %d not found", resourceLink)
	}

	if strings.Trim(name, " ") == "" {
		return errors.New("Name must not be empty")
	}

	attributes := fmt.Sprintf("{ \"name\": \"%s\" }", name)

	reqBody := strings.NewReader(attributes)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/resourcelinks/%d", h.baseURL, resourceLink), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// SetResourceLinkDescription sets the description for the specified Phillips Hue
// resource link
func (h *Connection) SetResourceLinkDescription(resourceLink int, description string) error {
	// Error checking
	if !h.doesResourceLinkExist(resourceLink) {
		return fmt.Errorf("Resource link %d not found", resourceLink)
	}

	if strings.Trim(description, " ") == "" {
		return errors.New("Description must not be empty")
	}

	attributes := fmt.Sprintf("{ \"description\": \"%s\" }", description)

	reqBody := strings.NewReader(attributes)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/resourcelinks/%d", h.baseURL, resourceLink), reqBody)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteResourceLink deletes a Phillips Hue resource link
func (h *Connection) DeleteResourceLink(resourceLink int) error {
	// Error checking
	if !h.doesResourceLinkExist(resourceLink) {
		return fmt.Errorf("Resource link %d not found", resourceLink)
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/resourcelinks/%d", h.baseURL, resourceLink), nil)
	if err != nil {
		return err
	}

	err = h.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (h *Connection) doesResourceLinkExist(resourceLink int) bool {
	// If GetResourceLink returns an error, then the resource link doesn't exist
	_, err := h.GetResourceLink(resourceLink)
	if err != nil {
		return false
	}

	return true
}

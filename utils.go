package hue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func (h *Connection) get(url string) ([]byte, error) {
	err := h.initializeHue()
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", h.baseURL, url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = h.checkForErrors(body); err != nil {
		return nil, err
	}

	return body, nil
}

func (h *Connection) execute(req *http.Request) error {
	err := h.initializeHue()
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return h.checkForErrors(body)
}

func (h *Connection) checkForErrors(data []byte) error {
	// Check for errors
	errMsg := ""
	errorMsgs := make([]map[string]map[string]interface{}, 10)

	err := json.Unmarshal(data, &errorMsgs)
	if err != nil {
		return nil
	}

	for _, e := range errorMsgs {
		if e["error"]["description"] != nil && e["error"]["description"] != "" {
			errMsg += fmt.Sprintf("%s\n", e["error"]["description"].(string))
		}
	}

	if errMsg != "" {
		return fmt.Errorf("Error: %s", errMsg)
	}

	return nil
}

// formatSlice formats an int slice as a JSON string
func (h *Connection) formatSlice(sli []int) string {
	str := "["
	for _, s := range sli {
		str += fmt.Sprintf("\"%d\",", s)
	}
	str = str[:len(str)-1]
	str += "]"

	return str
}

// formatStruct formats a struct as a JSON string
func (h *Connection) formatStruct(data interface{}) string {
	d := reflect.ValueOf(data)
	t := d.Type()
	str := ""

	if d.Kind() == reflect.Slice {
		str = "["
		for i := 0; i < d.Len(); i++ {
			str += fmt.Sprintf("%s,", h.formatStruct(d.Field(i)))
		}

		if len(str) > 1 {
			str = str[:len(str)-1]
		}
		str += "]"
	} else {
		str = "{"

		for i := 0; i < d.NumField(); i++ {
			str += fmt.Sprintf("\"%s\": ", strings.ToLower(t.Field(i).Name))

			switch d.Field(i).Kind() {
			case reflect.String:
				str += fmt.Sprintf("\"%s\",", d.Field(i).Interface())
			// Check for slice and array
			case reflect.Struct:
				str += fmt.Sprintf("%s,", h.formatStruct(d.Field(i).Interface()))
			default:
				str += fmt.Sprintf("%v,", d.Field(i).Interface())
			}
		}

		if len(str) > 1 {
			str = str[:len(str)-1]
		}
		str += "}"
	}

	return str
}

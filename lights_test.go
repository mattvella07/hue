package hue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type hueLightTestData struct {
	One Light `json:"1"`
}

type newLightTest struct {
	Name string `json:"name"`
}

type hueNewLightTestData struct {
	Five     newLightTest `json:"5"`
	LastScan string       `json:"lastscan"`
}

func createTestConnection(scenario int) (Connection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/lights":
			if scenario == 1 {
				// One light
				data := hueLightTestData{
					One: Light{
						State: hueLightState{
							On:        false,
							Bri:       100,
							Hue:       200,
							Sat:       300,
							Effect:    "",
							XY:        []float32{0.45},
							CT:        400,
							Alert:     "",
							ColorMode: "",
							Mode:      "",
							Reachable: true,
						},
						SWUpdate: hueLightSWUpdate{
							State:       "noupdates",
							LastInstall: "2018-06-04T6:14:11",
						},
						Type:             "Extended color",
						Name:             "Hue color lamp 1",
						ModelID:          "LCT016",
						ManufacturerName: "Phillips",
						ProductName:      "Hue color lamp",
						Capabilities: hueLightCapabilities{
							Certified: true,
							Control: hueLightCapabilitiesControl{
								MindimLevel:    1000,
								MaxLumen:       800,
								ColorGamutType: "C",
								ColorGamut:     [][]float32{[]float32{0.2, 0.3}, []float32{0.4, 0.5}},
								CT: hueLightCapabilitiesCT{
									Min: 153,
									Max: 500,
								},
							},
							Streaming: hueLightCapabilitiesStreaming{
								Renderer: true,
								Proxy:    true,
							},
						},
						Config: hueLightConfig{
							ArcheType: "sultanbulb",
							Function:  "mixed",
							Direction: "omnidirectional",
						},
						UniqueID:   "ab:cd:ef",
						SWVersion:  "1.29",
						SWConfigID: "ABCD",
						ProductID:  "Phillips-LCT016",
					},
				}

				returnData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("ERR: ", err)
				}

				w.Write(returnData)
			} else if scenario == 2 {
				// No lights
				w.Write(nil)
			}
		case "/lights/1", "/lights/2":
			if r.Method == "GET" {
				if scenario == 1 {
					// One light
					data := Light{
						State: hueLightState{
							On:        false,
							Bri:       100,
							Hue:       200,
							Sat:       300,
							Effect:    "",
							XY:        []float32{0.45},
							CT:        400,
							Alert:     "",
							ColorMode: "",
							Mode:      "",
							Reachable: true,
						},
						SWUpdate: hueLightSWUpdate{
							State:       "noupdates",
							LastInstall: "2018-06-04T6:14:11",
						},
						Type:             "Extended color",
						Name:             "Hue color lamp 1",
						ModelID:          "LCT016",
						ManufacturerName: "Phillips",
						ProductName:      "Hue color lamp",
						Capabilities: hueLightCapabilities{
							Certified: true,
							Control: hueLightCapabilitiesControl{
								MindimLevel:    1000,
								MaxLumen:       800,
								ColorGamutType: "C",
								ColorGamut:     [][]float32{[]float32{0.2, 0.3}, []float32{0.4, 0.5}},
								CT: hueLightCapabilitiesCT{
									Min: 153,
									Max: 500,
								},
							},
							Streaming: hueLightCapabilitiesStreaming{
								Renderer: true,
								Proxy:    true,
							},
						},
						Config: hueLightConfig{
							ArcheType: "sultanbulb",
							Function:  "mixed",
							Direction: "omnidirectional",
						},
						UniqueID:   "ab:cd:ef",
						SWVersion:  "1.29",
						SWConfigID: "ABCD",
						ProductID:  "Phillips-LCT016",
					}

					returnData, err := json.Marshal(data)
					if err != nil {
						fmt.Println("ERR: ", err)
					}

					w.Write(returnData)
				} else if scenario == 2 {
					// No lights
					w.Write(nil)
				}
			} else if r.Method == "PUT" {
				w.Write([]byte("[{\"success\":{\"/lights/1/name\":\"Bedroom Light\"}}]"))
			}
		case "/lights/new":
			// One new light
			if scenario == 1 {
				data := hueNewLightTestData{
					Five: newLightTest{
						Name: "Hue lamp 5",
					},
					LastScan: "2018-10-12T12:00:00",
				}

				returnData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("ERR: ", err)
				}

				w.Write(returnData)
			} else if scenario == 2 {
				// No new lights
				w.Write(nil)
			}
		}
	}))
	return Connection{
		UserID:            "TEST",
		internalIPAddress: "localhost",
		baseURL:           server.URL,
	}, server
}

func TestGetAllLights(t *testing.T) {
	t.Run("One light found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		lights, err := h.GetAllLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(lights) != expected {
				t.Fatalf("Expected %d light, got %d", expected, len(lights))
			}
		}

		{
			expected := "Hue color lamp 1"
			if lights[0].Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, lights[0].Name)
			}
		}

		{
			expected := "Extended color"
			if lights[0].Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, lights[0].Type)
			}
		}

		{
			expected := "Phillips"
			if lights[0].ManufacturerName != expected {
				t.Fatalf("Expected ManufacturerName to equal %s, got %s", expected, lights[0].ManufacturerName)
			}
		}
	})

	t.Run("No lights found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		lights, err := h.GetAllLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 0
			if len(lights) != expected {
				t.Fatalf("Expected %d lights, got %d", expected, len(lights))
			}
		}
	})
}

func TestGetLight(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		light, err := h.GetLight(1)
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := "Hue color lamp 1"
			if light.Name != expected {
				t.Fatalf("Expected Name to equal %s, got %s", expected, light.Name)
			}
		}

		{
			expected := "Extended color"
			if light.Type != expected {
				t.Fatalf("Expected Type to equal %s, got %s", expected, light.Type)
			}
		}

		{
			expected := "Phillips"
			if light.ManufacturerName != expected {
				t.Fatalf("Expected ManufacturerName to equal %s, got %s", expected, light.ManufacturerName)
			}
		}
	})

	t.Run("Not found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetLight(1)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestGetNewLights(t *testing.T) {
	t.Run("New light found", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		newLights, err := h.GetNewLights()
		if err != nil {
			t.Fatal(err)
		}

		{
			expected := 1
			if len(newLights.NewLights) != expected {
				t.Fatalf("Expected number of new lights to equal %d, got %d", expected, len(newLights.NewLights))
			}
		}

		{
			expected := 5
			if newLights.NewLights[0].ID != expected {
				t.Fatalf("Expected ID of first new light to equal %d, got %d", expected, newLights.NewLights[0].ID)
			}
		}

		{
			expected := "Hue lamp 5"
			if newLights.NewLights[0].Name != expected {
				t.Fatalf("Expected Name of first new light to equal %s, got %s", expected, newLights.NewLights[0].Name)
			}
		}

		{
			expected := "2018-10-12T12:00:00"
			if newLights.LastScan != expected {
				t.Fatalf("Expected LastScan to equal %s, got %s", expected, newLights.LastScan)
			}
		}
	})

	t.Run("No new lights found", func(t *testing.T) {
		h, server := createTestConnection(2)
		defer server.Close()

		_, err := h.GetNewLights()
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "No new lights found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestFindNewLights(t *testing.T) {
	h, server := createTestConnection(3)
	defer server.Close()

	err := h.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOnLight(t *testing.T) {
	t.Run("Light exists", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLight(3)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOnLightWithColor(t *testing.T) {
	t.Run("Light exists", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 200, 233)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(3, 0.3, 0.2, 100, 200, 233)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid x value", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 2, 0.2, 100, 200, 233)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Invalid color value: x must be between 0 and 1"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid y value", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 0.2, 3, 100, 200, 233)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Invalid color value: y must be between 0 and 1"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid bri value", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 300, 200, 233)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Invalid brightness value: bri must be between 1 and 254"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid hue value", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 65539, 233)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Invalid hue value: hue must be between 0 and 65,535"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Invalid sat value", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 200, 350)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Invalid saturation value: sat must be between 0 and 254"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestTurnOffLight(t *testing.T) {
	t.Run("Light exists", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOffLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.TurnOffLight(3)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestRenameLight(t *testing.T) {
	t.Run("Successful rename", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.RenameLight(1, "Light Renamed")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.RenameLight(3, "Light Renamed")
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})

	t.Run("Inavlid name", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.RenameLight(1, "")
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Name must not be empty"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

func TestDeleteLight(t *testing.T) {
	t.Run("Successful delete", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.DeleteLight(1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Light doesn't exist", func(t *testing.T) {
		h, server := createTestConnection(1)
		defer server.Close()

		err := h.DeleteLight(3)
		if err == nil {
			t.Fatal(err)
		}

		{
			expected := "Light not found"
			if err.Error() != expected {
				t.Fatalf("Expected error message to equal %s, got %s", expected, err.Error())
			}
		}
	})
}

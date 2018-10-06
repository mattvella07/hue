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

func createTestConnection(scenario int) (Connection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/lights":
			if scenario == 1 {
				//One light
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
				//No lights
				w.Write(nil)
			}
		case "/lights/1", "/lights/2":
			if scenario == 1 {
				//One light
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
				//No lights
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

func TestFindNewLights(t *testing.T) {
	h, server := createTestConnection(1)
	defer server.Close()

	err := h.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOnLight(t *testing.T) {
	//Update once error checking is added

	h, server := createTestConnection(1)
	defer server.Close()

	err := h.TurnOnLight(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOnLightWithColor(t *testing.T) {
	//Update once error checking is added

	h, server := createTestConnection(1)
	defer server.Close()

	err := h.TurnOnLightWithColor(1, 0.3, 0.2, 100, 200, 300)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTurnOffLight(t *testing.T) {
	//Update once error checking is added

	h, server := createTestConnection(1)
	defer server.Close()

	err := h.TurnOffLight(1)
	if err != nil {
		t.Fatal(err)
	}
}

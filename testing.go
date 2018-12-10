package hue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type lightTestData struct {
	One Light `json:"1"`
}

type newLightTest struct {
	Name string `json:"name"`
}

type newLightTestData struct {
	Five     newLightTest `json:"5"`
	LastScan string       `json:"lastscan"`
}

type groupTestData struct {
	One Group `json:"1"`
}

func createTestConnection(scenario int) (Connection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/lights":
			if scenario == 1 {
				// One light
				data := lightTestData{
					One: Light{
						State: lightState{
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
						SWUpdate: lightSWUpdate{
							State:       "noupdates",
							LastInstall: "2018-06-04T6:14:11",
						},
						Type:             "Extended color",
						Name:             "Hue color lamp 1",
						ModelID:          "LCT016",
						ManufacturerName: "Phillips",
						ProductName:      "Hue color lamp",
						Capabilities: lightCapabilities{
							Certified: true,
							Control: lightCapabilitiesControl{
								MindimLevel:    1000,
								MaxLumen:       800,
								ColorGamutType: "C",
								ColorGamut:     [][]float32{[]float32{0.2, 0.3}, []float32{0.4, 0.5}},
								CT: lightCapabilitiesCT{
									Min: 153,
									Max: 500,
								},
							},
							Streaming: lightCapabilitiesStreaming{
								Renderer: true,
								Proxy:    true,
							},
						},
						Config: lightConfig{
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
						State: lightState{
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
						SWUpdate: lightSWUpdate{
							State:       "noupdates",
							LastInstall: "2018-06-04T6:14:11",
						},
						Type:             "Extended color",
						Name:             "Hue color lamp 1",
						ModelID:          "LCT016",
						ManufacturerName: "Phillips",
						ProductName:      "Hue color lamp",
						Capabilities: lightCapabilities{
							Certified: true,
							Control: lightCapabilitiesControl{
								MindimLevel:    1000,
								MaxLumen:       800,
								ColorGamutType: "C",
								ColorGamut:     [][]float32{[]float32{0.2, 0.3}, []float32{0.4, 0.5}},
								CT: lightCapabilitiesCT{
									Min: 153,
									Max: 500,
								},
							},
							Streaming: lightCapabilitiesStreaming{
								Renderer: true,
								Proxy:    true,
							},
						},
						Config: lightConfig{
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
				data := newLightTestData{
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
		case "/groups":
			if r.Method == "GET" {
				if scenario == 1 {
					// One group
					data := groupTestData{
						One: Group{
							Name:   "Group 1",
							Lights: []string{"1", "2"},
							Type:   "LightGroup",
							Action: groupAction{
								On:        false,
								Bri:       100,
								Hue:       200,
								Sat:       250,
								Effect:    "none",
								XY:        []float32{0.3, 0.4},
								CT:        250,
								Alert:     "select",
								ColorMode: "ct",
							},
						},
					}

					returnData, err := json.Marshal(data)
					if err != nil {
						fmt.Println("ERR: ", err)
					}

					w.Write(returnData)
				} else if scenario == 2 {
					// No groups
					w.Write(nil)
				}
			} else if r.Method == "POST" {
				w.Write([]byte("[{\"success\":{\"id\":\"1\"}}]"))
			}
		case "/groups/1":
			if r.Method == "GET" {
				if scenario == 1 {
					// One Group
					data := Group{
						Name:   "Group 1",
						Lights: []string{"1", "2"},
						Type:   "LightGroup",
						Action: groupAction{
							On:        false,
							Bri:       100,
							Hue:       200,
							Sat:       250,
							Effect:    "none",
							XY:        []float32{0.3, 0.4},
							CT:        250,
							Alert:     "select",
							ColorMode: "ct",
						},
					}

					returnData, err := json.Marshal(data)
					if err != nil {
						fmt.Println("ERR: ", err)
					}

					w.Write(returnData)
				}
			} else if r.Method == "PUT" || r.Method == "DELETE" {
				w.Write([]byte("[{\"success\":{\"/groups/1/\":\"Success\"}}]"))
			}
		case "/groups/1/action":
			w.Write([]byte("[{\"success\":{\"/groups/1/action\":\"Success\"}}]"))
		}
	}))
	return Connection{
		UserID:            "TEST",
		internalIPAddress: "localhost",
		baseURL:           server.URL,
	}, server
}

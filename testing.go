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

type scheduleTestData struct {
	One Schedule `json:"1"`
}

type sceneTestData struct {
	One Scene `json:"1"`
}

type sensorTestData struct {
	One Sensor `json:"1"`
}

type newSensorTest struct {
	Name string `json:"name"`
}

type newSensorTestData struct {
	Five     newSensorTest `json:"5"`
	LastScan string        `json:"lastscan"`
}

type ruleTestData struct {
	One Rule `json:"1"`
}

type configurationWhitelistTestData struct {
	ABC ConfigurationWhitelist `json:"abc"`
}

type resourceLinkTestData struct {
	One ResourceLink `json:"1"`
}

func createTestConnection(scenario int) (Connection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if scenario == 1 {
				// Generate test data based on URL
				data := generateTestData(r.URL.String())
				if data == nil {
					w.Write(nil)
					return
				}

				returnData, err := json.Marshal(data)
				if err != nil {
					fmt.Println("ERR: ", err)
				}

				w.Write(returnData)
			} else if scenario == 2 {
				// No data returned from GET
				w.Write(nil)
			}
		case "PUT", "POST", "DELETE":
			if scenario == 1 {
				//Successful PUT, POST, or DELETE
				w.Write([]byte(fmt.Sprintf("[{\"success\":\"%s %s\"}]", r.Method, r.URL.String())))
			} else if scenario == 2 {
				//Error returned from PUT, POST, or DELETE
				w.Write([]byte(fmt.Sprintf("[{\"error\": {\"description\": \"Error while performing %s on %s\"}}]", r.Method, r.URL.String())))
			}
		}
	}))

	return Connection{
		UserID:            "TEST",
		internalIPAddress: "localhost",
		baseURL:           server.URL,
		isInitialized:     true,
	}, server
}

func generateTestData(url string) interface{} {
	switch url {
	case "/lights":
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

		return data
	case "/lights/1", "/lights/2":
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

		return data
	case "/lights/new":
		data := newLightTestData{
			Five: newLightTest{
				Name: "Hue lamp 5",
			},
			LastScan: "2018-10-12T12:00:00",
		}

		return data
	case "/groups":
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

		return data
	case "/groups/1":
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

		return data
	case "/schedules":
		data := scheduleTestData{
			One: Schedule{
				Name:        "Timer",
				Description: "Simple timer",
				Command: ScheduleCommand{
					Address: "/api/abc/groups/0/action",
					Body: ScheduleCommandBody{
						Scene: "1234",
					},
					Method: "PUT",
				},
				Time:       "PT00:01:00",
				Created:    "2018-12-10T13:39:16",
				Status:     "enabled",
				AutoDelete: false,
				StartTime:  "2018-12-10T14:00:00",
			},
		}

		return data
	case "/schedules/1":
		data := Schedule{
			Name:        "Timer",
			Description: "Simple timer",
			Command: ScheduleCommand{
				Address: "/api/abc/groups/0/action",
				Body: ScheduleCommandBody{
					Scene: "1234",
				},
				Method: "PUT",
			},
			Time:       "PT00:01:00",
			Created:    "2018-12-10T13:39:16",
			Status:     "enabled",
			AutoDelete: false,
			StartTime:  "2018-12-10T14:00:00",
		}

		return data
	case "/scenes":
		data := sceneTestData{
			One: Scene{
				Name:    "Night time",
				Type:    "LightScene",
				Group:   "1",
				Lights:  []string{"1", "2"},
				Owner:   "abcd",
				Recycle: true,
				Locked:  false,
				AppData: SceneAppData{
					Version: 1,
					Data:    "myAppData",
				},
				Picture:     "",
				LastUpdated: "2018-12-21T12:00:00",
				Version:     1,
			},
		}

		return data
	case "/scenes/1":
		data := Scene{
			Name:    "Night time",
			Type:    "LightScene",
			Group:   "1",
			Lights:  []string{"1", "2"},
			Owner:   "abcd",
			Recycle: true,
			Locked:  false,
			AppData: SceneAppData{
				Version: 1,
				Data:    "myAppData",
			},
			Picture:     "",
			LastUpdated: "2018-12-21T12:00:00",
			Version:     1,
		}

		return data
	case "/sensors":
		data := sensorTestData{
			One: Sensor{
				State: SensorState{
					Daylight:    true,
					LastUpdated: "2014-06-27T07:38:51",
				},
				Config: SensorConfig{
					On:            true,
					Long:          "none",
					Lat:           "none",
					SunriseOffset: 50,
					SunsetOffset:  50,
				},
				Name:             "Daylight",
				Type:             "Daylight",
				ModelID:          "PHDL00",
				ManufacturerName: "Phillips",
				SWVersion:        "1.0",
				ID:               1,
			},
		}

		return data
	case "/sensors/new":
		data := newSensorTestData{
			Five: newSensorTest{
				Name: "Sensor 5",
			},
			LastScan: "2018-10-12T12:00:00",
		}

		return data
	case "/sensors/1":
		data := Sensor{
			State: SensorState{
				Daylight:    true,
				LastUpdated: "2014-06-27T07:38:51",
			},
			Config: SensorConfig{
				On:            true,
				Long:          "none",
				Lat:           "none",
				SunriseOffset: 50,
				SunsetOffset:  50,
			},
			Name:             "Daylight",
			Type:             "Daylight",
			ModelID:          "PHDL00",
			ManufacturerName: "Phillips",
			SWVersion:        "1.0",
			ID:               1,
		}

		return data
	case "/rules":
		data := ruleTestData{
			One: Rule{
				Name:           "Rule 1",
				LastTriggered:  "2014-08-27T07:38:51",
				CreationTime:   "2014-06-27T07:38:51",
				TimesTriggered: 10,
				Owner:          "abc",
				Status:         "enabled",
				Conditions:     []RuleConditions{},
				Actions:        []RuleActions{},
				ID:             1,
			},
		}

		return data
	case "/rules/1":
		data := Rule{
			Name:           "Rule 1",
			LastTriggered:  "2014-08-27T07:38:51",
			CreationTime:   "2014-06-27T07:38:51",
			TimesTriggered: 10,
			Owner:          "abc",
			Status:         "enabled",
			Conditions:     []RuleConditions{},
			Actions:        []RuleActions{},
			ID:             1,
		}

		return data
	case "/config":
		data := Configuration{
			Name:          "Phillips hue",
			ZigbeeChannel: 15,
			Mac:           "abcd:efgh",
			DHCP:          true,
			IPAddress:     "192.0.0.1",
			NetMask:       "0.0.0.1",
			Gateway:       "0.0.0.1",
			ProxyAddress:  "none",
			ProxyPort:     0,
			UTC:           "2018-07-17T09:27:35",
			LocalTime:     "2018-07-17T09:27:35",
			Timezone:      "Central",
			SWVersion:     "100",
			APIVersion:    "1.3.0",
			SWUpdate: ConfigurationSWUpdate{
				UpdateState: 0,
				URL:         "",
				Text:        "",
				Notify:      true,
			},
			LinkButton:       true,
			PortalServices:   true,
			PortalConnection: "",
			PortalState: ConfigurationPortalState{
				SignedOn:      true,
				Incoming:      true,
				Outgoing:      true,
				Communication: "",
			},
		}

		return data
	case "/resourcelinks":
		data := resourceLinkTestData{
			One: ResourceLink{
				Name:        "Sunrise",
				Description: "Wake up experience",
				Class:       1,
				Owner:       "abcd",
				Links:       []string{"/schedules/1", "/schedules/2"},
			},
		}

		return data
	case "/resourcelinks/1":
		data := ResourceLink{
			Name:        "Sunset",
			Description: "Go to sleep experience",
			Class:       1,
			Owner:       "abcd",
			Links:       []string{"/schedules/1", "/schedules/2"},
		}

		return data
	}

	return nil
}

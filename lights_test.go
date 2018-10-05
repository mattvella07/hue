package hue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type hueLightTestData struct {
	One hueLight `json:"1"`
}

func createTestConnection() (Connection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/lights":
			data := hueLightTestData{
				One: hueLight{
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
		default:
			fmt.Println("No match")
		}
	}))
	return Connection{
		UserID:            "TEST",
		internalIPAddress: "localhost",
		baseURL:           server.URL,
	}, server
}

func TestGetAllLights(t *testing.T) {
	h, server := createTestConnection()
	defer server.Close()

	lights, err := h.GetAllLights()
	if err != nil {
		t.Fatal(err)
	}

  {
    if len(lights) != 1 {
      t.Fatalf("Expected 1 light, got %d", len(lights))
    }
  }

	{
		if lights[0].Name != "Hue color lamp 1" {
			t.Fatalf("Expected %s to equal Hue color lamp 1", lights[0].Name) 
		} 
	}
}

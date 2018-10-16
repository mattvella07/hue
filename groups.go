package hue

type groupAction struct {
	On        bool      `json:"on"`
	Bri       int       `json:"bri"`
	Hue       int       `json:"hue"`
	Sat       int       `json:"sat"`
	Effect    string    `json:"effect"`
	XY        []float32 `json:"xy"`
	CT        int       `json:"ct"`
	Alert     string    `json:"alert"`
	ColorMode string    `json:"colormode"`
}

// Group contains all data returned from the Phillips Hue API
// for an individual Phillips Hue light group
type Group struct {
	Name   string      `json:"name"`
	Lights []string    `json:"lights"`
	Type   string      `json:"type"`
	Action groupAction `json:"action"`
}

// GetAllGroups gets all Phillips Hue light groups connected to current bridge
func (h *Connection) GetAllGroups() ([]Group, error) {
	// GET - /api/<username>/groups

	return nil, nil
}

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

type Group struct {
	Name   string      `json:"name"`
	Lights []string    `json:"lights"`
	Type   string      `json:"type"`
	Action groupAction `json:"action"`
}

func (h *Connection) GetAllGroups() {
	// GET - /api/<username>/groups
}

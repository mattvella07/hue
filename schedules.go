package hue

type scheduleCommandBody struct {
	Scene string `json:"scene"`
}

type scheduleCommand struct {
	Address string              `json:"address"`
	Body    scheduleCommandBody `json:"body"`
	Method  string              `json:"method"`
}

type Schedule struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Command     scheduleCommand `json:"command"`
	Time        string          `json:"time"`
	Created     string          `json:"created"`
	Status      string          `json:"status"`
	AutoDelete  bool            `json:"autodelete"`
	StartTime   string          `json:"starttime"`
}

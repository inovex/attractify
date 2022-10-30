package responses

type CheckedAction struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name      string `json:"name"`
	UserValue string `json:"userValue"`
	DataValue string `json:"dataValue"`
	Blocking  bool   `json:"blocking"`
	Info      string `json:"info"`
}

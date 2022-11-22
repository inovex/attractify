package responses

type CheckedAction struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Display bool   `json:"display"`
	Steps   []Step `json:"steps"`
}

type Step struct {
	Name string `json:"name"`
	//UserValue string `json:"userValue"`
	//DataValue string `json:"dataValue"`
	State string `json:"state"`
	Info  string `json:"info"`
}

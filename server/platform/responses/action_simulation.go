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
	Blocking bool   `json:"blocking"`
	Info     string `json:"info"`
}

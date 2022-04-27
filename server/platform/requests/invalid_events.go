package requests

type InvalidEventUpdate struct {
	NewName string `json:"name" binding:"omitempty"`
	EventId string `json:"eventId" binding:"omitempty"`
}

package requests

type Data struct {
	Data interface{} `json:"data" binding:"required"`
}

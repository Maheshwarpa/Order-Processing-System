package Event

type OrderCreatedResponse struct {
	Order_id string `json:"order_id"`
	Status   string `json:"status"`
}

var OCRList []OrderCreatedResponse

type ProcessingResponse struct {
	Order_id       string `json: "order_id"`
	Payment_status string `json:"payment_status"`
}

var PRList []ProcessingResponse

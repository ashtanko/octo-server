package model

type IncomingRequest struct {
	Source        string  `json:"source"`
	State         string  `json:"state"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transactionId"`
}

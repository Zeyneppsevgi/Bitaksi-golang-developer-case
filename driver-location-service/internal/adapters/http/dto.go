package http

type pointDTO struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type batchItemDTO struct {
	DriverID string   `json:"driverId"`
	Location pointDTO `json:"location"`
}

type batchRequest struct {
	Items []batchItemDTO `json:"items"`
}

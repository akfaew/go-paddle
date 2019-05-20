package paddle

type Payment struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

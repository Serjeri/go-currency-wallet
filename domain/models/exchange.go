package models

type Exchange struct {
	Amount       int    `json:"amount"`
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
}

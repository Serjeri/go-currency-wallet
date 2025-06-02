package models


type UpdateBalance struct{
	Amount float64 `form:"amount" json:"amount" binding:"required"`
	Currency string `form:"currency" json:"currency" binding:"required"`
}

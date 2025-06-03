package models

type UpdateBalance struct {
	Amount   int    `form:"amount" json:"amount" binding:"required"`
	Currency string  `form:"currency" json:"currency" binding:"required"`
	Status   string  `form:"status" json:"status" binding:"required"`
}

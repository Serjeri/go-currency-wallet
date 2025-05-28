package models

type Login struct {
	Name     string `form:"name" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

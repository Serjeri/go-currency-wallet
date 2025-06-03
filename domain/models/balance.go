package models

type Balance struct{
	USD int `db:"usd"`
    RUB int `db:"rub"`
    EUR int `db:"eur"`
}

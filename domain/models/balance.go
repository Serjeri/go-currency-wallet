package models

type Balance struct{
	USD float64 `db:"usd"`
    RUB float64 `db:"rub"`
    EUR float64 `db:"eur"`
}

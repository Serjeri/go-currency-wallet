package lib

import (
	"errors"
	"fmt"
	"gw-currency-wallet/domain/models"
)


func DeductFromBalance(balance *models.Balance, currency string, amount int) (int, error) {
	var newAmount int

	switch currency {
	case "USD":
		if balance.USD < amount {
			return 0, errors.New("insufficient USD funds")
		}
		newAmount = balance.USD - amount
		balance.USD = newAmount
	case "RUB":
		if balance.RUB < amount {
			return 0, errors.New("insufficient RUB funds")
		}
		newAmount = balance.RUB - amount
		balance.RUB = newAmount
	case "EUR":
		if balance.EUR < amount {
			return 0, errors.New("insufficient EUR funds")
		}
		newAmount = balance.EUR - amount
		balance.EUR = newAmount
	default:
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}
	return newAmount, nil
}


func AddToBalance(amount int, currency string, balance *models.Balance ) int {
    switch currency {
    case "USD":
        balance.USD += amount
        return balance.USD
    case "RUB":
        balance.RUB += amount
        return balance.RUB
    case "EUR":
        balance.EUR += amount
        return balance.EUR
    default:
        return 0
    }
}


func UpdateBalance(currentBalance *models.Balance, updateBalance *models.UpdateBalance) (int, error) {
	var newAmount int

	switch updateBalance.Status {
	case "deposit":
		newAmount = AddToBalance(updateBalance.Amount, updateBalance.Currency, currentBalance)
	case "withdrawal":
		newAmount, _ = DeductFromBalance(currentBalance, updateBalance.Currency, updateBalance.Amount)
	default:
		return 0, fmt.Errorf("unknown status: %s", updateBalance.Status)
	}
	return newAmount, nil
}

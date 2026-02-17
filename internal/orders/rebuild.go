package orders

import "errors"

func Rebuild(id string, amount int64, status Status) (*Order, error) {
	if amount < 0 {
		return nil, errors.New("invalid amount")
	}

	return &Order{
		id:     id,
		amount: amount,
		status: status,
	}, nil
}

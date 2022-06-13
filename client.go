package easebuzz

import (
	"github.com/gagan1147/easebuzz-go/payment"
	"github.com/gagan1147/easebuzz-go/va"
)

func NewVAClient(key, salt string) *va.Client {
	vaClient := &va.Client{
		Key:  key,
		Salt: salt,
	}
	return vaClient
}

func NewPaymentClient(key, salt string) *payment.Client {
	paymentClient := &payment.Client{
		Key:  key,
		Salt: salt,
	}
	return paymentClient
}

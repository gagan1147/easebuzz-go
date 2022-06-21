package easebuzz

import (
	"github.com/gagan1147/easebuzz-go/payment"
	"github.com/gagan1147/easebuzz-go/va"
)

// NewVAClient returns a pointer to the VA client
func NewVAClient(key, salt string) *va.Client {
	vaClient := &va.Client{
		Key:  key,
		Salt: salt,
	}
	return vaClient
}

// NewPaymentClient returns a pointer to the Payment client
func NewPaymentClient(key, salt string) *payment.Client {
	paymentClient := &payment.Client{
		Key:  key,
		Salt: salt,
	}
	return paymentClient
}

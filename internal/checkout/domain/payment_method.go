package domain

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "CreditCard"
	PaymentMethodBankTransfer PaymentMethod = "BankTransfer"
)

func (m PaymentMethod) String() string {
	return string(m)
}

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCreditCard, PaymentMethodBankTransfer:
		return true
	}

	return false
}

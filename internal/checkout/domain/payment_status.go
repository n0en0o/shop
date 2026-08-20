package domain

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "Pending"
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed    PaymentStatus = "Failed"
	PaymentStatusRefunded  PaymentStatus = "Refunded"
)

func (s PaymentStatus) String() string {
	return string(s)
}

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending,
		PaymentStatusCompleted,
		PaymentStatusFailed,
		PaymentStatusRefunded:
		return true
	}

	return false
}

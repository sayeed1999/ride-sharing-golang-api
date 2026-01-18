package trip

type TripStatus int

const (
	ONGOING TripStatus = iota + 1
	WAITING_FOR_PAYMENT
	PAYMENT_COMPLETED
)

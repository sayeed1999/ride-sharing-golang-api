package trip

import "fmt"

type ITransitionChecker interface {
	IsTransitionValid(fromStatus TripStatus, toStatus TripStatus) bool
}

type TripStatusTransitionChecker struct {
	tripMap map[TripStatus][]TripStatus
}

func TripStatusTransitionCheckerMap() *TripStatusTransitionChecker {
	return &TripStatusTransitionChecker{
		tripMap: map[TripStatus][]TripStatus{
			ONGOING: {
				WAITING_FOR_PAYMENT,
			},
			WAITING_FOR_PAYMENT: {
				PAYMENT_COMPLETED,
			},
			PAYMENT_COMPLETED: {},
		},
	}
}

func (checker *TripStatusTransitionChecker) IsTransitionValid(fromStatus TripStatus, toStatus TripStatus) (bool, error) {
	supportedStatuses, ok := checker.tripMap[fromStatus]
	if !ok {
		return false, fmt.Errorf("transition from %s to %s is not supported", fromStatus.String(), toStatus.String())
	}

	for _, supported := range supportedStatuses {
		if supported == toStatus {
			return true, nil
		}
	}

	return false, fmt.Errorf("transition from %s to %s is not allowed", fromStatus.String(), toStatus.String())
}

func (status TripStatus) String() string {
	statusStrings := [...]string{
		"ONGOING",
		"WAITING_FOR_PAYMENT",
		"PAYMENT_COMPLETED",
	}
	if status < ONGOING || status > PAYMENT_COMPLETED {
		return "Unknown"
	}
	return statusStrings[status-1]
}

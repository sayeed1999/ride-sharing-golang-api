package trip_request

import "fmt"

type TripRequestStatusTransitionChecker struct {
	tripRequestMap map[TripRequestStatus][]TripRequestStatus
}

func TripRequestStatusTransitionCheckerMap() *TripRequestStatusTransitionChecker {
	return &TripRequestStatusTransitionChecker{
		tripRequestMap: map[TripRequestStatus][]TripRequestStatus{
			NO_DRIVER_FOUND: {
				CUSTOMER_CANCELED,
				DRIVER_ACCEPTED,
			},
			CUSTOMER_CANCELED: {},
			DRIVER_ACCEPTED: {
				CUSTOMER_REJECTED_DRIVER,
				DRIVER_REJECTED_CUSTOMER,
			},
			CUSTOMER_REJECTED_DRIVER: {},
			DRIVER_REJECTED_CUSTOMER: {
				NO_DRIVER_FOUND,
			},
			TRIP_STARTED:          {},
			TRIP_REQUEST_REJECTED: {},
		},
	}
}

func (checker *TripRequestStatusTransitionChecker) IsTransitionValid(fromStatus, toStatus TripRequestStatus) (bool, error) {
	supportedStatuses, ok := checker.tripRequestMap[fromStatus]
	if !ok {
		return false, fmt.Errorf("transition from %s to %s is not supported", fromStatus.String(), toStatus.String())
	}

	for _, status := range supportedStatuses {
		if status == toStatus {
			return true, nil
		}
	}

	return false, fmt.Errorf("transition from %s to %s is not allowed", fromStatus.String(), toStatus.String())
}

func (status TripRequestStatus) String() string {
	statusStrings := [...]string{
		"NO_DRIVER_FOUND",
		"CUSTOMER_CANCELED",
		"DRIVER_ACCEPTED",
		"CUSTOMER_REJECTED_DRIVER",
		"DRIVER_REJECTED_CUSTOMER",
		"TRIP_STARTED",
		"TRIP_REQUEST_REJECTED",
	}
	if status < NO_DRIVER_FOUND || status > TRIP_REQUEST_REJECTED {
		return "Unknown"
	}
	return statusStrings[status-1]
}

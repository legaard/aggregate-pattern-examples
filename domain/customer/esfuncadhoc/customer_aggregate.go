package esfuncadhoc

import (
	"aggregate_implementation_examples/domain/command"
	"aggregate_implementation_examples/domain/event"
	"aggregate_implementation_examples/stream"
)

// State changes to Aggregate as events (ES) done in a functional style where state is created ad-hoc
// inside each function
//
// State is partially build inside each function for idempotency and (maybe) new events. If the state change is
// valid new events are "returned" to be appended to the event stream
//
// Note: Creation of state will in some cases be duplicated - can be cumbersome and lead to errors due to state discrepancies

func RegisterCustomer(events []stream.Event, cmd command.RegisterCustomer) ([]interface{}, error) {
	// build state on the fly
	var alreadyRegistered bool
	for _, streamEvent := range events {
		if _, ok := streamEvent.Data.(event.CustomerRegistered); ok {
			alreadyRegistered = true
		}
	}
	// idempotency check
	if alreadyRegistered {
		return nil, nil
	}

	// return new events to append
	return []interface{}{
		event.CustomerRegistered{
			CustomerID:       cmd.CustomerID,
			FullName:         cmd.FullName,
			EmailAddress:     cmd.EmailAddress,
			ConfirmationHash: cmd.ConfirmationHash,
		},
	}, nil
}

func ConfirmEmail(events []stream.Event, cmd command.ConfirmationEmailAddress) ([]interface{}, error) {
	// build state on the fly
	var (
		customerId       string
		alreadyConfirmed bool
		confirmationHash string
	)
	for _, streamEvent := range events {
		switch e := streamEvent.Data.(type) {
		case event.CustomerRegistered:
			alreadyConfirmed = false
			confirmationHash = e.ConfirmationHash
			customerId = e.CustomerID
		case event.CustomerEmailAddressChanged:
			alreadyConfirmed = false
			confirmationHash = e.ConfirmationHash
		case event.CustomerEmailAddressConfirmed:
			alreadyConfirmed = true
		}
	}
	// idempotency check
	if alreadyConfirmed && confirmationHash == cmd.ConfirmationHash {
		return nil, nil
	}

	// return new events to append
	if confirmationHash == cmd.ConfirmationHash {
		return []interface{}{
			event.CustomerEmailAddressConfirmed{
				CustomerID: customerId,
			},
		}, nil
	}

	return []interface{}{
		event.CustomerEmailAddressConfirmationFailed{
			CustomerID: customerId,
		},
	}, nil
}

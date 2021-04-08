package esfuncrecons

import (
	"aggregate_implementation_examples/domain/command"
	"aggregate_implementation_examples/domain/event"
	"aggregate_implementation_examples/stream"
)

// State changes to Aggregate as events (ES) done in a functional style where state is created via a shared (across functions)
// reconstitution mechanism
//
// State is fully build inside each function for idempotency and (maybe) new events. If the state change is
// valid new events are "returned" to be appended to the event stream

func RegisterCustomer(events []stream.Event, cmd command.RegisterCustomer) ([]interface{}, error) {
	// build state
	state := Reconstitute(events)

	// idempotency check
	if state.CustomerID != "" {
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
	// build state
	state := Reconstitute(events)

	// idempotency check
	if state.EmailConfirmed && state.EmailAddress == cmd.ConfirmationHash {
		return nil, nil
	}

	// return new events
	if state.ConfirmationHash == cmd.ConfirmationHash {
		return []interface{}{
			event.CustomerEmailAddressConfirmed{
				CustomerID: state.CustomerID,
			},
		}, nil
	}

	return []interface{}{
		event.CustomerEmailAddressConfirmationFailed{
			CustomerID: state.CustomerID,
		},
	}, nil
}

package esfuncext

import (
	"aggregate_implementation_examples/domain/command"
	"aggregate_implementation_examples/domain/event"
)

// State changes to Aggregate as events (ES) done in a functional style where state is constructed outside
// the function
//
// Aggregate state is derived (projection) from events and passed to a function which might "emit" new events
//
// Note: Idempotency is important to avoid duplicates of events, e.g. two event.CustomerRegistered on the same
// stream. This goes for all ES solutions

func RegisterCustomer(state State, cmd command.RegisterCustomer) ([]interface{}, error) {
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

func ConfirmEmail(state State, cmd command.ConfirmationEmailAddress) ([]interface{}, error) {
	// idempotency check
	if state.EmailConfirmed && state.ConfirmationHash == cmd.ConfirmationHash {
		return nil, nil
	}

	// return new events to append
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

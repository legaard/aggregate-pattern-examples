package esfuncrecons

import (
	"aggregate_implementation_examples/domain/event"
	"aggregate_implementation_examples/stream"
)

type State struct {
	CustomerID       string
	EmailAddress     string
	EmailConfirmed   bool
	ConfirmationHash string
}

func Reconstitute(events []stream.Event) State {
	var state State
	for _, streamEvent := range events {
		switch e := streamEvent.Data.(type) {
		case event.CustomerRegistered:
			state.CustomerID = e.CustomerID
			state.EmailConfirmed = false
			state.ConfirmationHash = e.ConfirmationHash
		case event.CustomerEmailAddressChanged:
			state.EmailConfirmed = false
			state.ConfirmationHash = e.ConfirmationHash
		case event.CustomerEmailAddressConfirmed:
			state.EmailConfirmed = true
		}
	}

	return state
}

package crudfunc

import (
	"aggregate_implementation_examples/domain/command"
	"fmt"
)

// State changes to Aggregate as CRUD done in a functional style
//
// Takes current state of the CustomerAggregate as argument, validates a state change is allowed
// and then returns a new state.

func RegisterCustomer(state State, cmd command.RegisterCustomer) (*State, error) {
	return &State{
		CustomerID:       cmd.CustomerID,
		EmailAddress:     cmd.EmailAddress,
		ConfirmationHash: cmd.ConfirmationHash,
	}, nil
}

func ConfirmEmail(state State, cmd command.ConfirmationEmailAddress) (*State, error) {
	if state.ConfirmationHash != cmd.ConfirmationHash {
		return nil, fmt.Errorf("could not confirm email %s using confirmation hash %q", state.EmailAddress, cmd.ConfirmationHash)
	}
	// return not state to be saved
	return &State{
		CustomerID:       state.CustomerID,
		EmailAddress:     state.EmailAddress,
		ConfirmationHash: cmd.ConfirmationHash,
		EmailConfirmed:   true,
	}, nil
}

package crudoop

import (
	"aggregate_implementation_examples/domain/command"
	"fmt"
)

// State changes to Aggregate as CRUD done in a OOP style
//
// Aggregate state is read into the struct from storage. Properties of the CustomerAggregate is then updated,
// if the command is valid, after which the Aggregate is saved again
//
// Note: If Integration Events must be emitted when a change occurs some way of marking the
// Aggregate as "dirty" (i.e. a state change has been performed) would probably be needed

type CustomerAggregate struct {
	CustomerID       string
	EmailAddress     string
	EmailConfirmed   bool
	ConfirmationHash string
}

func (c *CustomerAggregate) RegisterCustomer(cmd command.RegisterCustomer) error {
	c.CustomerID = cmd.CustomerID
	c.ConfirmationHash = cmd.ConfirmationHash
	c.EmailAddress = cmd.EmailAddress
	return nil
}

func (c *CustomerAggregate) ConfirmEmail(cmd command.ConfirmationEmailAddress) error {
	if c.ConfirmationHash != cmd.ConfirmationHash {
		return fmt.Errorf("could not confirm email %s using confirmation hash %q", c.EmailAddress, cmd.ConfirmationHash)
	}
	c.EmailConfirmed = true
	return nil
}

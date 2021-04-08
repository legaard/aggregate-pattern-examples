package esoop

import (
	"aggregate_implementation_examples/domain/command"
	"aggregate_implementation_examples/domain/event"
	"aggregate_implementation_examples/stream"
)

// State changes to Aggregate as events (ES) done in a OOP style where state is created via an Apply method on the
// Aggregate struct
//
// State is fully build via the Apply method and new events - in cases there is any - is returned via the methods
// handling the commands

type CustomerAggregate struct {
	customerID       string
	emailAddress     string
	emailConfirmed   bool
	confirmationHash string
	events           []interface{}
}

func (c *CustomerAggregate) RegisterCustomer(cmd command.RegisterCustomer) ([]interface{}, error) {
	// idempotency check
	if c.customerID != "" {
		return nil, nil
	}

	// new events to append
	return []interface{}{
		event.CustomerRegistered{
			CustomerID:       cmd.CustomerID,
			FullName:         cmd.FullName,
			EmailAddress:     cmd.EmailAddress,
			ConfirmationHash: cmd.ConfirmationHash,
		},
	}, nil
}

func (c *CustomerAggregate) ConfirmEmail(cmd command.ConfirmationEmailAddress) ([]interface{}, error) {
	// idempotency check
	if c.emailConfirmed && c.emailAddress == cmd.ConfirmationHash {
		return nil, nil
	}

	// new events to append
	if c.confirmationHash == cmd.ConfirmationHash {
		return []interface{}{
			event.CustomerEmailAddressConfirmed{
				CustomerID: c.customerID,
			},
		}, nil
	}

	return []interface{}{
		event.CustomerEmailAddressConfirmationFailed{
			CustomerID: c.customerID,
		},
	}, nil
}

func (c *CustomerAggregate) Apply(streamEvent stream.Event) {
	// build aggregate state
	switch e := streamEvent.Data.(type) {
	case event.CustomerRegistered:
		c.customerID = e.CustomerID
		c.emailConfirmed = false
		c.confirmationHash = e.ConfirmationHash
	case event.CustomerEmailAddressChanged:
		c.emailConfirmed = false
		c.confirmationHash = e.ConfirmationHash
	case event.CustomerEmailAddressConfirmed:
		c.emailConfirmed = true
	}
}

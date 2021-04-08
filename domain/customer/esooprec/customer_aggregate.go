package esooprec

import (
	"aggregate_implementation_examples/domain/command"
	"aggregate_implementation_examples/domain/event"
	"aggregate_implementation_examples/stream"
)

// State changes to Aggregate as events (ES) done in a OOP style where state is created via an Apply method on the
// Aggregate struct
//
// State is fully build via the Apply method and new events - in cases there is any - are recorded on the struct for
// external callers to retrieve via the Events method

type CustomerAggregate struct {
	customerID       string
	emailAddress     string
	emailConfirmed   bool
	confirmationHash string
	events           []interface{}
}

func (c *CustomerAggregate) RegisterCustomer(cmd command.RegisterCustomer) error {
	// idempotency check
	if c.customerID != "" {
		return nil
	}

	// new events to append
	c.Emit(
		event.CustomerRegistered{
			CustomerID:       cmd.CustomerID,
			FullName:         cmd.FullName,
			EmailAddress:     cmd.EmailAddress,
			ConfirmationHash: cmd.ConfirmationHash,
		})

	return nil
}

func (c *CustomerAggregate) ConfirmEmail(cmd command.ConfirmationEmailAddress) error {
	// idempotency check
	if c.emailConfirmed && c.emailAddress == cmd.ConfirmationHash {
		return nil
	}

	// new events to append
	if c.confirmationHash == cmd.ConfirmationHash {
		c.Emit(event.CustomerEmailAddressConfirmed{
			CustomerID: c.customerID,
		})
	} else {
		c.Emit(event.CustomerEmailAddressConfirmationFailed{
			CustomerID: c.customerID,
		})
	}

	return nil
}

func (c *CustomerAggregate) Emit(events ...interface{}) {
	for _, streamEvent := range events {
		c.events = append(c.events, streamEvent)
	}
}

// Accessed externally
func (c *CustomerAggregate) Events() []interface{} {
	return c.events
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

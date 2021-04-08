package event

type CustomerEmailAddressChanged struct {
	CustomerID       string
	EmailAddress     string
	ConfirmationHash string
}

package event

type CustomerRegistered struct {
	CustomerID       string
	FullName         string
	EmailAddress     string
	ConfirmationHash string
}

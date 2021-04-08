package command

type RegisterCustomer struct {
	CustomerID       string
	FullName         string
	EmailAddress     string
	ConfirmationHash string
}

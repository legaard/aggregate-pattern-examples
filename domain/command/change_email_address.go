package command

type ChangeEmailAddress struct {
	CustomerID       string
	EmailAddress     string
	ConfirmationHash string
}

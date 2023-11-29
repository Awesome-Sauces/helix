package transaction

type Transaction struct {
	sender   string
	receiver string

	instructions string
}

/*

A transaction is not a set payment, the sender of a transaction
must provide the subsequent `instructions` for their transaction.
The sender may choose to include a payment mechanism in the transaction
instructions or a complex entanglement of funds between multiple accounts.

*/

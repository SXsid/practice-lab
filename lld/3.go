package main

// Composition
type BankAccount struct {
	balance float64
}
type Customer struct {
	name     string
	accounts []*BankAccount // holds actual objects
}

func (c *Customer) OpenAccount(balance float64) *BankAccount {
	account := &BankAccount{balance: balance} // created inside
	c.accounts = append(c.accounts, account)
	return account
}

// Aggregation
func (c *Customer) AddAccount(account *BankAccount) { // passed from outside
	c.accounts = append(c.accounts, account)
}

package entity

type User struct {
	Id       int
	Username string
	Password string
}

type Received struct {
	FromUser string `json:"from"`
	Amount   int    `json:"amount"`
}

type Send struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransactionHistory struct {
	Received []Received `json:"received"`
	Send     []Send     `json:"sent"`
}

type Info struct {
	Balance     int                `json:"balance"`
	Inventory   []Inventory        `json:"inventory"`
	CoinHistory TransactionHistory `json:"coinHistory"`
}

type Inventory struct {
	ItemName string `json:"name"`
	Quantity int    `json:"quantity"`
}

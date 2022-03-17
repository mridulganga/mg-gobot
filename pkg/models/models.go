package models

import "time"

type User struct {
	ID        string
	Name      string
	Balance   Balance
	Last      Last
	Inventory Inventory
}

type Balance struct {
	Wallet int
	Bank   int
}

type Last struct {
	Beg     time.Time
	Search  time.Time
	Gamble  time.Time
	Daily   time.Time
	Weekly  time.Time
	Monthly time.Time
}

type Inventory map[string]struct {
	Quantity int
	Price    int
}

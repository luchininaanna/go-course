package model

type Order struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
}

type MenuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderList struct {
	Orders []Order `json:"orders"`
}

type DetailedOrder struct {
	Order
	Time int `json:"orderedAtTimestamp"`
	Cost int `json:"cost"`
}

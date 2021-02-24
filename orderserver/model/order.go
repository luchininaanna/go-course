package model

type OrderListItem struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
}

type OrderItem struct {
	ID               string     `json:"id"`
	OrderAtTimestamp int        `json:"orderAtTimestamp"`
	Cost             int        `json:"cost"`
	MenuItems        []MenuItem `json:"menuItems"`
}

type MenuItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

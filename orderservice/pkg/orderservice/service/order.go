package service

type Order struct {
	MenuItems []MenuItem
}

type MenuItem struct {
	ID       string
	Quantity int
}

type OrderList struct {
	Orders []Order `json:"orders"`
}

type DetailedOrder struct {
	Order Order `json:"order"`
	Time  int   `json:"orderedAtTimestamp"`
	Cost  int   `json:"cost"`
}

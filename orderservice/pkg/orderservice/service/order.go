package service

import "time"

type Order struct {
	ID        string     `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
	Time      time.Time  `json:"orderedAtTimestamp"`
	Cost      int        `json:"cost"`
}

type MenuItem struct {
	ID       string
	Quantity int
}

type OrderList struct {
	Orders []Order `json:"orders"`
}

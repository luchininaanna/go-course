package model

import (
	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	MenuItems []MenuItem
}

type MenuItem struct {
	ID       uuid.UUID
	Quantity int
}

type OrderRepository interface {
	AddOrder(order Order) error
}

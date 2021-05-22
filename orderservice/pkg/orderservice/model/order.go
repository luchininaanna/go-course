package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID        uuid.UUID
	MenuItems []MenuItem
	Time      time.Time
	Cost      int
}

type MenuItem struct {
	ID       uuid.UUID
	Quantity int
}

type OrderRepository interface {
	AddOrder(order Order) error
	UpdateOrder(order Order) error
	DeleteOrder(uuid uuid.UUID) error
	GetOrder(uuid uuid.UUID) (*Order, error)
	GetOrders() ([]Order, error)
}

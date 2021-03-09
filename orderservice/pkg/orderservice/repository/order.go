package repository

import (
	"database/sql"
	"orderserver/pkg/orderservice/model"
	"time"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) model.OrderRepository {
	return &orderRepository{db: db}
}

func (o orderRepository) AddOrder(order model.Order) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	orderIdBin, err := order.ID.MarshalBinary()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = o.db.Exec("INSERT INTO `order`(id, cost, created_at, updated_at, deleted_at) VALUES (?, 77, ?, null, null);", orderIdBin, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, orderItem := range order.MenuItems {
		orderItemIdBin, err := orderItem.ID.MarshalBinary()
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = o.db.Exec("INSERT INTO `order_item`(order_id, menu_item_id, quantity) VALUES (?, ?, ?);", orderIdBin, orderItemIdBin, orderItem.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

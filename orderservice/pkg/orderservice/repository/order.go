package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"orderserver/pkg/orderservice/model"
	"strconv"
	"strings"
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

func (o orderRepository) UpdateOrder(order model.Order) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	orderIdBin, err := order.ID.MarshalBinary()
	if err != nil {
		tx.Rollback()
		return err
	}

	//удалить все позиции по заказу
	_, err = o.db.Exec("DELETE FROM `order_item` WHERE order_id = ?;", orderIdBin)
	if err != nil {
		tx.Rollback()
		return err
	}

	//добавить позиции по заказу
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

func (o orderRepository) DeleteOrder(orderUuid uuid.UUID) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	orderIdBin, err := orderUuid.MarshalBinary()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = o.db.Exec("DELETE FROM `order` WHERE id = ?;", orderIdBin)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (o orderRepository) GetOrder(orderUuid uuid.UUID) (*model.Order, error) {
	orderIdBin, err := orderUuid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query(""+
		"SELECT "+
		"BIN_TO_UUID(o.id) AS id, "+
		"GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.menu_item_id), \"=\", oi.quantity)) AS menuItems, "+
		"o.created_at AS time, "+
		"o.cost AS cost "+
		"FROM `order` o "+
		"LEFT JOIN order_item oi ON o.id = oi.order_id "+
		"WHERE o.id = ? "+
		"GROUP BY o.id", orderIdBin)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return nil, err
}

func (o orderRepository) GetOrders() ([]model.Order, error) {
	rows, err := o.db.Query("" +
		"SELECT " +
		"BIN_TO_UUID(o.id) AS id, " +
		"GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.menu_item_id), \"=\", oi.quantity)) AS menuItems, " +
		"o.created_at AS time, " +
		"o.cost AS cost " +
		"FROM `order` o " +
		"LEFT JOIN order_item oi ON o.id = oi.order_id " +
		"GROUP BY o.id")

	if err != nil {
		return nil, err
	}

	return parseOrders(rows)
}

func parseOrders(r *sql.Rows) ([]model.Order, error) {
	var orders []model.Order

	for r.Next() {
		order, err := parseOrder(r)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func parseOrder(r *sql.Rows) (*model.Order, error) {
	var orderId string
	var menuItems string
	var t time.Time
	var cost int

	err := r.Scan(&orderId, &menuItems, &t, &cost)
	if err != nil {
		return nil, err
	}

	orderUuid, err := uuid.Parse(orderId)
	if err != nil {
		return nil, err
	}

	menuItemsArray := strings.Split(menuItems, ",")

	var modelMenuItems []model.MenuItem
	for _, menuItem := range menuItemsArray {
		s := strings.Split(menuItem, "=")
		itemUuid, err := uuid.Parse(s[0])
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		modelMenuItems = append(modelMenuItems, model.MenuItem{ID: itemUuid, Quantity: quantity})
	}

	return &model.Order{
		ID:        orderUuid,
		MenuItems: modelMenuItems,
		Time:      t,
		Cost:      cost,
	}, nil
}

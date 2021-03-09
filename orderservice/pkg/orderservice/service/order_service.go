package service

import (
	"errors"
	"github.com/google/uuid"
	"orderserver/pkg/orderservice/model"
)

type orderService struct {
	repository model.OrderRepository
}

type OrderService interface {
	AddOrder(orderData Order) error
}

func NewOrderService(r model.OrderRepository) OrderService {
	return &orderService{repository: r}
}

func (os *orderService) AddOrder(orderData Order) error {
	err := validateOrder(orderData)
	if err != nil {
		return err
	}

	orderModel, err := createOrderModel(orderData)
	if err != nil {
		return err
	}

	return os.repository.AddOrder(*orderModel)
}

func createOrderModel(orderData Order) (*model.Order, error) {
	menuItems := make([]model.MenuItem, 0)
	for _, menuItem := range orderData.MenuItems {
		itemUuid, err := uuid.Parse(menuItem.ID)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, model.MenuItem{ID: itemUuid, Quantity: menuItem.Quantity})
	}

	orderUuid := uuid.New()
	return &model.Order{ID: orderUuid, MenuItems: menuItems}, nil
}

func validateOrder(orderData Order) error {
	err := isValidItemsQuantity(orderData)
	if err != nil {
		return err
	}

	err = isOrderHasNotDuplicatedItems(orderData)
	if err != nil {
		return err
	}

	return nil
}

func isValidItemsQuantity(orderData Order) error {
	for _, menuItem := range orderData.MenuItems {
		if menuItem.Quantity == 0 {
			return errors.New("Order has item with zero quantity")
		}
	}

	return nil
}

func isOrderHasNotDuplicatedItems(orderData Order) error {
	items := make(map[string]bool)
	for _, entry := range orderData.MenuItems {
		if _, value := items[entry.ID]; !value {
			items[entry.ID] = true
		} else {
			return errors.New("Order has duplicated items")
		}
	}

	return nil
}

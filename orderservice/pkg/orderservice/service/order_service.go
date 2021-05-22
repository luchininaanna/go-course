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
	UpdateOrder(id string, orderData Order) error
	DeleteOrder(id string) error
	GetOrderInfo(id string) (*Order, error)
	GetOrders() ([]Order, error)
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

func (os *orderService) UpdateOrder(id string, orderData Order) error {
	orderUuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	err = validateOrder(orderData)
	if err != nil {
		return err
	}

	orderModel, err := updateOrderModel(orderUuid, orderData)
	if err != nil {
		return err
	}

	return os.repository.UpdateOrder(*orderModel)
}

func (os *orderService) DeleteOrder(id string) error {
	orderUuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return os.repository.DeleteOrder(orderUuid)
}

func (os *orderService) GetOrderInfo(id string) (*Order, error) {
	orderUuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	order, err := os.repository.GetOrder(orderUuid)
	if err != nil {
		return nil, err
	}

	var menuItems []MenuItem
	for _, item := range order.MenuItems {
		menuItems = append(menuItems, MenuItem{
			item.ID.String(),
			item.Quantity,
		})
	}

	return &Order{
		order.ID.String(),
		menuItems,
		order.Time,
		order.Cost,
	}, err
}

func (os *orderService) GetOrders() ([]Order, error) {
	var or []Order
	orders, err := os.repository.GetOrders()
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		var menuItems []MenuItem
		for _, item := range order.MenuItems {
			menuItems = append(menuItems, MenuItem{
				item.ID.String(),
				item.Quantity,
			})
		}
		o := Order{
			order.ID.String(),
			menuItems,
			order.Time,
			order.Cost,
		}
		or = append(or, o)
	}

	return or, nil
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

func updateOrderModel(orderUuid uuid.UUID, orderData Order) (*model.Order, error) {
	menuItems := make([]model.MenuItem, 0)
	for _, menuItem := range orderData.MenuItems {
		itemUuid, err := uuid.Parse(menuItem.ID)
		if err != nil {
			return nil, err
		}
		menuItems = append(menuItems, model.MenuItem{ID: itemUuid, Quantity: menuItem.Quantity})
	}

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

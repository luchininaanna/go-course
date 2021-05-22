package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"orderserver/pkg/orderservice/service"
	"testing"
)

type mockOrderService struct{}

func (m mockOrderService) GetOrders() ([]service.Order, error) {
	return []service.Order{
		{
			ID: uuid.New().String(),
			MenuItems: []service.MenuItem{
				{ID: uuid.New().String(), Quantity: 26},
				{ID: uuid.New().String(), Quantity: 5},
			},
		},
	}, nil
}

func (m mockOrderService) AddOrder(orderData service.Order) error {
	panic("implement me")
}

func (m mockOrderService) UpdateOrder(id string, orderData service.Order) error {
	panic("implement me")
}

func (m mockOrderService) DeleteOrder(id string) error {
	panic("implement me")
}

func (m mockOrderService) GetOrderInfo(id string) (*service.Order, error) {
	panic("implement me")
}

func TestOrderList(t *testing.T) {
	srv := server{orderService: mockOrderService{}}
	w := httptest.NewRecorder()
	srv.getOrders(w, nil)

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var items []service.Order
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}

package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"orderserver/pkg/orderservice/service"
	"testing"
)

func TestOrderList(t *testing.T) {
	w := httptest.NewRecorder()
	getOrders(w, nil)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have: %d, want: %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	items := service.OrderList{}
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}

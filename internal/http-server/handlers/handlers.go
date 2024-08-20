package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"wb-order-service/internal/storage/memcache"
)

//type OrderTest struct {
//	Id    int     `json:"id"`
//	Name  string  `json:"name"`
//	Price float64 `json:"price"`
//}

//type Response struct {
//	Order  OrderTest `json:"order"`
//	Status int       `json:"status"`
//	Error  string    `json:"error,omitempty"`
//}

func OrderPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/order.html")
}

func GetOrderHandler(log *slog.Logger, c memcache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		order, err := c.GetOrderById(id)
		if err != nil {
			log.Error("err", err, "Failed to get order by id: ", id)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		log.Info("Got order by id: ", id)

		render.JSON(w, r, order)
	}
}

//func GetOrderHandler(log *slog.Logger) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		idParam := chi.URLParam(r, "id")
//		id, err := strconv.Atoi(idParam)
//		if err != nil {
//			log.Error("failed to parse id", err)
//			render.JSON(w, r, Response{
//				Status: http.StatusInternalServerError,
//				Error:  "internal server error",
//			})
//			return
//		}
//
//		order, err := GetOrder(id)
//		if err != nil {
//			log.Error("failed to get order", err)
//			render.JSON(w, r, Response{
//				Status: http.StatusInternalServerError,
//				Error:  "internal server error",
//			})
//			return
//		}
//
//		log.Info("received")
//		render.JSON(w, r, Response{
//			Order:  order,
//			Status: http.StatusOK,
//		})
//	}
//}

//// TODO: update method
//func GetOrder(id int) (OrderTest, error) {
//	return OrderTest{
//		Id:    id,
//		Name:  "OrderName",
//		Price: 25.5,
//	}, nil
//}

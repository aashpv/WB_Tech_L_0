package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"wb-order-service/internal/storage/memcache"
)

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

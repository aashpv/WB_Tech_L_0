package memcache

import (
	"fmt"
	gocache "github.com/patrickmn/go-cache"
	"strconv"
	"time"
	"wb-order-service/internal/model"
	"wb-order-service/internal/storage"
	"wb-order-service/internal/storage/postgres"
)

type Cache interface {
	SaveOrder(order model.Order) error
	GetOrderById(id string) (model.Order, error)
	LoadOrders() error
}

type cache struct {
	ch  *gocache.Cache
	pgs *postgres.Database
}

func NewCache(pgs *postgres.Database) Cache {
	return &cache{
		ch:  gocache.New(5*time.Minute, 10*time.Minute),
		pgs: pgs,
	}
}

func (c *cache) SaveOrder(order model.Order) error {
	id, err := c.pgs.InsertOrderToDB(&order)
	if err != nil {
		return fmt.Errorf("error saving order to db: %v", err)
	}
	c.ch.Set(strconv.Itoa(id), order, gocache.DefaultExpiration)

	return nil
}

func (c *cache) GetOrderById(id string) (model.Order, error) {
	data, found := c.ch.Get(id)
	if !found {
		return model.Order{}, storage.ErrOrderNotFound
	}
	return data.(model.Order), nil
}

func (c *cache) LoadOrders() error {
	orders, err := c.pgs.GetOrdersFromDB()
	if err != nil {
		return fmt.Errorf("error getting orders: %v", err)
	}

	for i, order := range orders {
		c.ch.Set(strconv.Itoa(i), order, gocache.DefaultExpiration)
	}

	return nil
}

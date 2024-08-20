package stan

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"wb-order-service/internal/model"
	"wb-order-service/internal/storage/memcache"
)

type Nats interface {
	SaveOrder(order model.Order) error
}

type Stan struct {
	ClusterId string
	ClientId  string
	Subject   string
}

func (stn *Stan) NewStan(ch memcache.Cache) error {
	sc, err := stan.Connect(stn.ClusterId, stn.ClientId)
	if err != nil {
		return fmt.Errorf("error connecting to NATS Streaming: %s", err)
	} else {
		fmt.Println("Connected to NATS Streaming")
	}

	var order model.Order
	_, err = sc.Subscribe(stn.Subject, func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Printf("Error unmarshalling: %s", err)
			return
		}
		fmt.Println("Unmarshalling successfully")

		if err := validator.New().Struct(&order); err != nil {
			fmt.Printf("Error validating: %s", err)
			return
		}
		fmt.Println("Validating successfully")

		err = ch.SaveOrder(order)
		if err != nil {
			fmt.Printf("Error saving: %s", err)
			return
		}
		fmt.Println("Saving successfully")
	})
	if err != nil {
		return fmt.Errorf("error subscribing: %s", err)
	} else {
		fmt.Println("Subscribed to NATS Streaming")
	}

	return nil
}

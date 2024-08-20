package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

type Order struct {
	OrderUid          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery" validate:"required"`
	Payment           Payment  `json:"payment" validate:"required"`
	Items             []Item   `json:"items" validate:"required"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	ShardKey          string   `json:"shardkey" validate:"required"`
	SmId              int      `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"required"`
	PaymentDt    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required"`
	GoodsTotal   int    `json:"goods_total" validate:"required"`
	CustomFee    int    `json:"custom_fee" validate:"required"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"min=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"min=0"`
	NmId        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

const (
	url        = "nats://localhost:4222"
	cluster_id = "test-cluster"
	client_id  = "wb-order-test-service"
	subject    = "wb-tech"
)

func main() {
	sc, err := stan.Connect(cluster_id, client_id, stan.NatsURL(url))
	if err != nil {
		log.Printf("error connecting: %v", err)
	} else {
		log.Println("connected successfully")
	}
	defer sc.Close()

	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Printf("error reading file: %v", err)
	} else {
		log.Println("reading successfully")
	}

	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		log.Printf("error parsing json: %v", err)
	} else {
		log.Println("parsing successfully")
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Printf("error marshalling json: %v", err)
	} else {
		log.Println("marshalling successfully")
	}

	err = sc.Publish(subject, jsonData)
	if err != nil {
		log.Printf("error publishing message: %v", err)
	} else {
		log.Println("published message")
	}
}

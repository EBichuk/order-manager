package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	producer *kafka.Writer
}

func NewProducer() *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:29092", "localhost:39092", "localhost:19092"},
		Topic:   "order",
	})

	return &Producer{
		producer: w,
	}
}

func (p *Producer) Produce(message *Order) error {
	msg, err := json.Marshal(message)

	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Value: []byte(msg),
	}

	err = p.producer.WriteMessages(context.Background(), kafkaMsg)

	if err != nil {
		log.Fatal("Ошибка при отправке:", err)
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}

func MakeRandomOrder() Order {
	item := Item{
		ChrtID:      1000000 + rand.IntN(100000),
		TrackNumber: "WBTESTTRACK",
		Price:       1 + rand.IntN(1000),
		Rid:         uuid.New().String(),
		NameItem:    "Test Name Item",
		Sale:        rand.IntN(100),
		Size:        0,
		TotalPrice:  rand.IntN(10000),
		NmID:        0,
		Brand:       "Test Brand",
		Status:      202,
	}
	payment := Payment{
		Transaction:  uuid.New().String(),
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       rand.IntN(100),
		PaymentDt:    1000000 + rand.IntN(100000),
		Bank:         "test bank",
		DeliveryCost: rand.IntN(1000),
		GoodsTotal:   rand.IntN(10000),
		CustomFee:    0,
	}
	delivery := Delivery{
		Name:    "Test Testov",
		Phone:   "+79000000000",
		Zip:     "123456",
		City:    "Moscow",
		Address: "Lenina 10",
		Region:  "Moscow",
		Email:   "test-testov@gmail.com",
	}
	order := Order{
		OrderUID:          uuid.New().String(),
		TrackNumber:       "new",
		Entry:             "WBIL",
		Locate:            "en",
		InternalSignature: " ",
		CustomerID:        uuid.New().String(),
		DeliveryService:   "meest",
		Shardkey:          "0",
		SmID:              rand.IntN(100),
		DateCreated:       time.Now().UTC(),
		OffShard:          "1",
		Delivery:          delivery,
		Payment:           payment,
		Item:              []Item{item},
	}
	return order
}

func main() {
	p := NewProducer()
	for i := 0; i < 20; i++ {
		order := MakeRandomOrder()

		time.Sleep(5 * time.Second)

		err := p.Produce(&order)
		if err != nil {
			return
		}
	}
	p.Close()
}

type Delivery struct {
	ID       int    `json:"-"`
	OrderUID string `json:"-"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Zip      string `json:"zip" validate:"required"`
	City     string `json:"city" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Region   string `json:"region" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string `json:"-"`
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int    `json:"amount" validate:"required"`
	PaymentDt    int    `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   int    `json:"goods_total" validate:"required"`
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Item              []Item    `json:"items" validate:"required,min=1"`
	Locate            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OffShard          string    `json:"oof_shard" validate:"required"`
}

type Item struct {
	ID          int    `json:"-"`
	OrderUID    string `json:"-"`
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	NameItem    string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"required,gte=0,lte=100"`
	Size        int    `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required"`
	NmID        int    `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

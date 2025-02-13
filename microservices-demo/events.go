package events

import (
    "encoding/json"

    "github.com/streadway/amqp"
)

type OrderCreated struct {
    OrderID string `json:"order_id"`
    UserID  string `json:"user_id"`
    Amount  float64 `json:"amount"`
}

func PublishOrderCreated(order OrderCreated) error {
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        return err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "order_created",
        true,
        false,
        false,
        false,
        nil,
    )

    body, _ := json.Marshal(order)
    return ch.Publish(
        "",
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
}
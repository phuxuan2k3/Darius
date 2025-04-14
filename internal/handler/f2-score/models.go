package f2_score

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ScoreRequest struct {
	Msg amqp.Delivery
}

type ScoreResponse interface {
}

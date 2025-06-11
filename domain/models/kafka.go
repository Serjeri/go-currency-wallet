package models

import (
	"github.com/google/uuid"
	"time"
)

type KafkaEvent struct {
	EventID   uuid.UUID         `json:"event_id"`
	EventType string            `json:"event_type"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   KafkaEventPayload `json:"payload"`
}

type KafkaEventPayload struct {
	UserID       int    `json:"userId"`
	Amount       int    `json:"amount"`
	FromCurrency string `json:"FromCurrency"`
	ToCurrency   string `json:"ToCurrency"`
}

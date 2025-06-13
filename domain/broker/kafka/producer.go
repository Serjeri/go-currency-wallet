package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/IBM/sarama"
	"gw-currency-wallet/domain/models"
)

func Producer(userID int, exchange *models.Exchange) (bool, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"kafka:29092"}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	defer producer.Close()

	event := models.KafkaEvent{
		EventID:   uuid.New(),
		EventType: "exchange",
		Timestamp: time.Now().UTC(),
		Payload: models.KafkaEventPayload{
			UserID:       userID,
			Amount:       exchange.Amount,
			FromCurrency: exchange.FromCurrency,
			ToCurrency:   exchange.ToCurrency,
		},
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return false, fmt.Errorf("failed to marshal kafka event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "exchange.transaction",
		Value: sarama.StringEncoder(jsonData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return false,fmt.Errorf("failed to send message: %w", err)
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", msg.Topic, partition, offset)
	return true, nil
}

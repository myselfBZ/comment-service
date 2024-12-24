package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/myselfBZ/comment-service/internal/storage"
)

func InitCommentProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	borkerURL := []string{"localhost:9092"}
	conn, err := sarama.NewSyncProducer(borkerURL, config)
	if err != nil {
		log.Fatal("error connecting to kafka: ", err)
	}
	return conn

}

func NewCommentProducer(prod sarama.SyncProducer, topic string) *CommentProducer {
	return &CommentProducer{prod, topic}
}

type CommentProducer struct {
	prod  sarama.SyncProducer
	topic string
}

func (p *CommentProducer) Push(c *storage.Comment) error {
	byteData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(byteData),
	}

	partitio, offset, err := p.prod.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message stored! Partition: %v Offset: %v", partitio, offset)
	return nil
}

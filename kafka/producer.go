package kafka

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string
	Topic   string
}

type Producer struct {
	w     *kafka.Writer
	topic string
}

func NewProducer(cfg ProducerConfig) *Producer {

	return &Producer{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      cfg.Brokers,
			Topic:        cfg.Topic,
			BatchTimeout: time.Millisecond * 5,
		}),
		topic: cfg.Topic,
	}
}

func (p *Producer) Close() {
	p.w.Close()
}

func (p *Producer) ProduceEv(evName string, data []byte) error {
	var (
		msg kafka.Message
		ev  Event

		evB []byte

		err error
	)

	ev.Name = evName
	ev.Data = data
	ev.CreatedAt = time.Now().String()

	if evB, err = proto.Marshal(&ev); err != nil {
		return err
	}

	msg.Value = evB

	if err = p.w.WriteMessages(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

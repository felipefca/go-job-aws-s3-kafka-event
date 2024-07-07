package kafka

import (
	"bytes"
	"context"
	"fmt"
	"go-job-aws-s3-kafka-event/configs"
	"go-job-aws-s3-kafka-event/internal/appctx"
	"go-job-aws-s3-kafka-event/internal/infra/kafka/avro"
	"go-job-aws-s3-kafka-event/internal/models"

	"github.com/IBM/sarama"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
)

type KafkaProducer interface {
	PublishEvent(ctx context.Context, model models.Message) error
}

type kafkaProducer struct {
	producer       sarama.SyncProducer
	config         configs.Kafka
	schemaClient   schemaregistry.Client
	schemaMetadata schemaregistry.SchemaMetadata
}

func NewKafkaProducer(producer sarama.SyncProducer, config configs.Kafka) (KafkaProducer, error) {
	schemaRegistry := schemaregistry.NewConfig(config.SchemaRegistryUrl)

	schemaClient, err := schemaregistry.NewClient(schemaRegistry)
	if err != nil {
		return nil, err
	}

	schemaLimit, err := schemaClient.GetLatestSchemaMetadata(config.TopicName + "-value")
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{
		producer:       producer,
		config:         config,
		schemaClient:   schemaClient,
		schemaMetadata: schemaLimit,
	}, nil
}

func (k *kafkaProducer) PublishEvent(ctx context.Context, model models.Message) error {
	logger := appctx.FromContext(ctx)

	event := avro.NewEvent()
	event.Message = model.Message

	var buf bytes.Buffer
	if err := event.Serialize(&buf); err != nil {
		return err
	}
	binary := buf.Bytes()

	payload := &sarama.ProducerMessage{
		Topic: k.config.TopicName,
		Value: sarama.ByteEncoder(binary),
	}

	partition, offset, err := k.producer.SendMessage(payload)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Success to send message. Partition: %d - offset: %d", partition, offset))
	return nil
}

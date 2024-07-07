package main

import (
	"context"
	"go-job-aws-s3-kafka-event/configs"
	"go-job-aws-s3-kafka-event/internal/appctx"
	"go-job-aws-s3-kafka-event/internal/server"

	"github.com/IBM/sarama"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logConfig := zap.NewProductionConfig()

	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ctx = appctx.WithLogger(ctx, logger)

	s3Client, err := connectAwsS3(ctx, configs.GetConfig().AwsS3)
	if err != nil {
		panic(err)
	}

	logger.Info("AWS S3 Client connected!")

	kafkaClient, err := connectKafka()
	if err != nil {
		panic(err)
	}
	defer kafkaClient.Close()

	logger.Info("Kafka Client connected!")

	s := server.NewServer(server.ServerOptions{
		Logger:        logger,
		Context:       ctx,
		S3Client:      s3Client,
		KafkaProducer: kafkaClient,
	})
	s.Start()
}

func connectAwsS3(ctx context.Context, cfg configs.AwsS3) (*s3.Client, error) {
	provider := credentials.NewStaticCredentialsProvider(cfg.ClientId, cfg.SecretKey, "")

	awsCfg, err := aws_config.LoadDefaultConfig(ctx, aws_config.WithRegion(cfg.Region), aws_config.WithCredentialsProvider(provider))
	if err != nil {
		return nil, err
	}

	awsClient := s3.NewFromConfig(awsCfg)
	return awsClient, nil
}

func connectKafka() (sarama.SyncProducer, error) {
	cfg := configs.GetConfig().Kafka

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	brokers := []string{cfg.Broker}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	return producer, nil
}

package server

import (
	"context"
	"os"

	"go-job-aws-s3-kafka-event/configs"
	"go-job-aws-s3-kafka-event/internal/appctx"
	"go-job-aws-s3-kafka-event/internal/infra/aws"
	"go-job-aws-s3-kafka-event/internal/infra/kafka"
	"go-job-aws-s3-kafka-event/internal/services"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type Server interface {
	Start()
}

type ServerOptions struct {
	Logger        *zap.Logger
	Context       context.Context
	S3Client      *s3.Client
	KafkaProducer sarama.SyncProducer
}

type server struct {
	ServerOptions
	services.Integrator
}

func NewServer(opt ServerOptions) Server {
	cfg := configs.GetConfig()

	var kafkaproducer, err = kafka.NewKafkaProducer(opt.KafkaProducer, cfg.Kafka)
	if err != nil {
		panic(err)
	}

	var s3Storage = aws.NewStorage(opt.S3Client, cfg.AwsS3.BucketName, cfg.AwsS3.Folder)
	var integrator = services.NewIntegrator(opt.Context, s3Storage, kafkaproducer)

	return server{
		ServerOptions: opt,
		Integrator:    integrator,
	}
}

func (s server) Start() {
	logger := appctx.FromContext(s.ServerOptions.Context)

	logger.Info("Starting job...")

	err := s.Integrator.Run(s.Context)
	if err != nil {
		logger.Error(err.Error())
		s.Context.Err()
		os.Exit(1)
	} else {
		logger.Info("Job successfully executed!")
		s.Context.Done()
	}
}

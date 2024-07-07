package services

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"go-job-aws-s3-kafka-event/configs"
	"go-job-aws-s3-kafka-event/internal/appctx"
	"go-job-aws-s3-kafka-event/internal/infra/aws"
	"go-job-aws-s3-kafka-event/internal/infra/kafka"
	"go-job-aws-s3-kafka-event/internal/models"
	"io"
	"path/filepath"
	"strings"
)

type Integrator interface {
	Run(ctx context.Context) error
}

type integrator struct {
	storage       aws.Storage
	kafkaproducer kafka.KafkaProducer
}

func NewIntegrator(ctx context.Context, storage aws.Storage, kafkaproducer kafka.KafkaProducer) Integrator {
	return &integrator{
		storage:       storage,
		kafkaproducer: kafkaproducer,
	}
}

func (i integrator) Run(ctx context.Context) error {
	logger := appctx.FromContext(ctx)
	s3Config := configs.GetConfig().AwsS3

	logger.Info("Processing integration...")

	resp, err := i.storage.ReaderListObjects(ctx, s3Config.Folder)
	if err != nil {
		return err
	}

	for _, obj := range resp.Contents {
		fileName := filepath.Base(*obj.Key)
		if strings.Contains(*obj.Key, s3Config.DestinationFolder) || !strings.Contains(*obj.Key, ".csv") {
			continue
		}

		result, err := i.storage.Reader(ctx, *obj.Key)
		if err != nil {
			return err
		}

		reader := csv.NewReader(result)
		reader.FieldsPerRecord = -1
		reader.Comma = ';'

		logger.Info(fmt.Sprintf("Start reading file %s", fileName))

		err = i.processRecords(ctx, result)
		if err != nil {
			return err
		}

		err = i.storage.MoveFile(ctx, fileName, s3Config.DestinationFolder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *integrator) processRecords(ctx context.Context, file io.ReadCloser) error {
	logger := appctx.FromContext(ctx)

	scanner := bufio.NewScanner(file)
	isFirstLine := true

	//Read file
	for scanner.Scan() {
		line := scanner.Text()

		//Replace double quotes with empty strings
		line = strings.Replace(line, "\"\"", "", -1)

		//Skip header
		if isFirstLine {
			isFirstLine = false
			continue
		}

		fields := strings.Split(line, ";")
		message := i.BuildMessage(fields)

		i.kafkaproducer.PublishEvent(ctx, *message)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	logger.Info("Processing completed!")

	return nil
}

func (i *integrator) BuildMessage(fields []string) *models.Message {
	return &models.Message{
		Message: fields[0],
	}
}

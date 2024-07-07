package configs

import (
	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	Application Application
	AwsS3       AwsS3
	Kafka       Kafka
}

type Application struct {
	AppName string
}

type AwsS3 struct {
	BucketName        string
	Folder            string
	DestinationFolder string
	Region            string
	ClientId          string
	SecretKey         string
}

type Kafka struct {
	Broker            string
	SchemaRegistryUrl string
	GroupId           string
	TopicName         string
}

func GetConfig() config {
	return *cfg
}

func init() {
	viper.SetDefault("PORT", "8080")

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfg = &config{
		Application: Application{
			AppName: "go-worker-rabbitmq-kafka-event",
		},
		AwsS3: AwsS3{
			BucketName:        viper.GetString("BUCKET_NAME"),
			Folder:            viper.GetString("BUCKET_FOLDER"),
			DestinationFolder: viper.GetString("DESTINATION_FOLDER"),
			Region:            viper.GetString("AWS_REGION"),
			ClientId:          viper.GetString("AWS_CLIENT_ID"),
			SecretKey:         viper.GetString("AWS_SECRET_KEY"),
		},
		Kafka: Kafka{
			Broker:            viper.GetString("KAFKA_BROKER"),
			SchemaRegistryUrl: viper.GetString("KAFKA_SCHEMA_REGISTRY_URL"),
			GroupId:           viper.GetString("KAFKA_GROUPID"),
			TopicName:         viper.GetString("KAFKA_TOPIC_NAME"),
		},
	}
}

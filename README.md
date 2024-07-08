# go-job-aws-s3-kafka-event
A Go job to read files from AWS S3 and stream contents in Apache Kafka

[![LinkedIn][linkedin-shield]][linkedin-url]


<!-- ABOUT THE PROJECT -->
## About the Project

A Go job to read files from folders in AWS S3 bucket and stream the contents in Apache Kafka

### Related Projects
- https://github.com/felipefca/go-api-event

- https://github.com/felipefca/go-worker-rabbitmq-kafka-event

![Screenshot_3](https://github.com/felipefca/go-api-event/assets/21323326/691c6cfe-2bce-48bc-b437-b60964738db4)

### Using



* [![Go][Go-badge]][Go-url]
* [![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-v3.0-orange.svg)](https://kafka.apache.org/)
* [![AWS S3](https://img.shields.io/badge/AWS-S3-FF9900.svg)](https://aws.amazon.com/s3/)
* [![Docker][Docker-badge]][Docker-url]

<!-- GETTING STARTED -->
## Getting Started

Instructions for running the application

### Prerequisites

Run the command to initialize Kafka (zookeeper, schema registry...) and the application on the selected port
* docker
  ```sh
  docker-compose up -d
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/felipefca/go-job-aws-s3-kafka-event.git
   ```

2. Exec
   ```sh
   go mod tidy
   ```
   
   ```sh
   cp .env.example .env
   ```
      
   ```sh
   go run ./cmd/main.go
   ```


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/felipe-fernandes-fca/
[Go-url]: https://golang.org/
[Go-badge]: https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white
[RabbitMQ-badge]: https://img.shields.io/badge/rabbitmq-%23ff6600.svg?style=flat&logo=rabbitmq&logoColor=white
[RabbitMQ-url]: https://www.rabbitmq.com/
[Docker-badge]: https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/

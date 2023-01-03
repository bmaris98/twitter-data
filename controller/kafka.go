package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

type KafkaContext struct {
	Reader   *kafka.Reader
	Writer   *kafka.Writer
	MongoCtx *MongoContext
}

func (kc *KafkaContext) Init(mongoCtx MongoContext) {
	kc.CreateReader()
	kc.StartWriter()
}

func (kc *KafkaContext) CreateReader() {
	kc.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"broker:9092"},
		Topic:   "twitter-streams-out",
		GroupID: "twitter-consumer-group",
	})
}

func (kc *KafkaContext) StartReading() {
	log.Println("Starting reading")
	for {
		m, err := kc.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("error while reading messages", err)
		}
		log.Println("Read done")
		query := string(m.Key)
		rawValue := string(m.Value)
		splittedRaw := strings.Split(rawValue, "===")
		timestamp, err := strconv.ParseUint(splittedRaw[1], 10, 64)
		if err != nil {
			log.Fatal("Error parsing timestamp to uint", err)
		}
		value, err := strconv.ParseUint(splittedRaw[0], 10, 64)
		stat := Stat{Query: query, Timestamp: timestamp, Value: value}
		mongoCtx.InsertUnsafeStat(stat)
	}
}

func (kc *KafkaContext) CloseReader() {
	if err := kc.Reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func (kc *KafkaContext) StartWriter() {
	// make a writer that produces to topic-A, using the least-bytes distribution
	kc.Writer = &kafka.Writer{
		Addr:                   kafka.TCP("broker:9092"),
		Topic:                  "twitter-streams-in",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

func (kc *KafkaContext) CloseWriter() {
	if err := kc.Writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func (kc *KafkaContext) PushMsg(query string, data string) {
	err := kc.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(query),
			Value: []byte(data),
		},
	)
	if err != nil {
		log.Println("failed to write messages:", err)
	}
}

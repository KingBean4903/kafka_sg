package main

import (
	"context"
	"fmt"
//	"os"
	"time"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/hamba/avro/v2"
)


type VideoTranscoded struct {
	VideoID string `avro:"video_id"`
	MasterURL string `avro:"master_url"`
	TranscodedAt string `avro:"transcoded_at"`
}

func main() {

	schema, err := avro.Parse(`{
			"type": "record",
			"name": "VideoTranscoded",
			"fields" : [
					{"name" : "video_id", "type" : "string"},
					{"name" : "master_url", "type" : "string"},
					{"name" : "transcoded_at", "type" : "string"}
			]
		}`)

	if err != nil {
			log.Fatalf("Error parsing schema: %v", err)
	}

	event := VideoTranscoded{
			VideoID: "video-458",
			MasterURL: "https://cdn.example.com/vid-458/manifest.mpd",
			TranscodedAt: time.Now().Format(time.RFC3339),
	}

	encoded, err := avro.Marshal(schema, event)
	if err != nil {
		log.Fatalf("Failed to encode Avro: %v", err)
	}

	writer := kafka.Writer{
				Addr: kafka.TCP("kafka-zoo:9092"),
				Topic: "videos.events",
				Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
					Key:  []byte(event.VideoID),
					Value: encoded,
	})

	if err != nil {
			fmt.Printf("Failed to send message: %v", err)
	}

	fmt.Println("Produced VideoTranscoded event", event.VideoID)
}

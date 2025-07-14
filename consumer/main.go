package main

import (
	"context"
	"fmt"
	"log"
	"github.com/hamba/avro/v2"
	"github.com/segmentio/kafka-go"
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

	r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{"kafka-zoo:9092"},
				Topic: "videos.events",
				GroupID: "video-service",
	})

	defer r.Close()

	log.Println("Listening for messages ...")

	for {
			
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("kafka err: %v", err)
				return
			}
			var record VideoTranscoded
			err = avro.Unmarshal(schema, msg.Value, &record)

			if err != nil { 
				log.Printf("Avro decoding err %v", err)
				continue
			}

			fmt.Printf("Record: %+v\n", record)
		
	}




}

package main

import (
	"context"
	"fmt"
	"encoding/binary"
	"log"
	"github.com/hamba/avro/v2"
	"github.com/segmentio/kafka-go"
	"github.com/riferrei/srclient"
)

type VideoTranscoded struct {
	VideoID string `avro:"video_id"`
	MasterURL string `avro:"master_url"`
	TranscodedAt string `avro:"transcoded_at"`
}

func main() {

//	registry := srclient.NewSchemaRegistryClient("http://localhost:8081")
	registry := srclient.NewSchemaRegistryClient("http://schema-registry:8081")

	/* schema, err := avro.Parse(`{
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
	} */

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
			/*var record VideoTranscoded
			err = avro.Unmarshal(schema, msg.Value, &record)

			if err != nil { 
				log.Printf("Avro decoding err %v", err)
				continue
			}*/

			decoded, err := decodeConfluentAvro(msg.Value, registry)
			if err != nil {
					log.Printf("Decode error: %v", err)
			}

			if record, ok := decoded.(VideoTranscoded); ok {
					fmt.Printf("Record transcoded video: %+v\n", record)
			} else { 
				fmt.Printf("Unexpected record type: %T\n", record)
			}
	}
}

func decodeConfluentAvro(data []byte, registry *srclient.SchemaRegistryClient) (interface{}, error) { 

		if len(data) < 5 || data[0] != 0 {
				return nil, fmt.Errorf("Invalid Confluent Avro format")
		}

		schemaID := int(binary.BigEndian.Uint32(data[1:5]))

		schema, err := registry.GetSchema(schemaID)

		if err != nil {
				return nil, fmt.Errorf("Schemas failed fetch: %v", err)
		}

		codec, err := avro.Parse(schema.Schema())
		if err != nil {
				return nil, fmt.Errorf("Failed to compile schema: %v", err)
		}

		// Decode payload
		var record VideoTranscoded
		err = avro.Unmarshal(codec, data[5:], &record)
		if err != nil {
				return nil, fmt.Errorf("avro unmarshal failed: %v", err)
		}

		return record, nil

}

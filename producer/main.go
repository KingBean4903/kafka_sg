package main

import (
	"context"
	"fmt"
//	"os"
	"encoding/binary"
	"bytes"
	"time"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/hamba/avro/v2"
	"github.com/riferrei/srclient"
)


type VideoTranscoded struct {
	VideoID string `avro:"video_id"`
	MasterURL string `avro:"master_url"`
	TranscodedAt string `avro:"transcoded_at"`
}

func main() {

	schemaBytes := `{
			"type": "record",
			"name": "VideoTranscoded",
			"fields" : [
					{"name" : "video_id", "type" : "string"},
					{"name" : "master_url", "type" : "string"},
					{"name" : "transcoded_at", "type" : "string"}
			]
		}`


//	schemaStr := string(schemaBytes)

	topic := "videos.events"

	schemaRegistry := srclient.NewSchemaRegistryClient("http://schema-registry:8081")
	schema, err := schemaRegistry.GetLatestSchema(topic)
	
	if err != nil {
		schema, err = schemaRegistry.CreateSchema(topic, schemaBytes, srclient.Avro)
		if err != nil { 
				log.Fatalf("Could not register schema: %v", err)
		}
	}

	schemaID := schema.ID()

	// Compile Avro code
	codec, err := avro.Parse(schema.Schema())
	if err != nil {
			log.Fatalf("Could not compile Avro codec: %v", err)
	}

	event := VideoTranscoded{
			VideoID: "video-458",
			MasterURL: "https://cdn.example.com/vid-458/manifest.mpd",
			TranscodedAt: time.Now().Format(time.RFC3339),
	}

	// Serialize with avro
	avroBytes, err := avro.Marshal(codec, event)

	writer := kafka.Writer{
				Addr: kafka.TCP("kafka-zoo:9092"),
				Topic: "videos.events",
				Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	// Wire format
	var buf bytes.Buffer
	buf.WriteByte(0)
	_ = binary.Write(&buf, binary.BigEndian, int32(schemaID))
	buf.Write(avroBytes)


	err = writer.WriteMessages(context.Background(), kafka.Message{
					Key:  []byte(event.VideoID),
					Value: buf.Bytes(),
	})

	if err != nil {
			fmt.Printf("Failed to send message: %v", err)
	}

	fmt.Println("Produced VideoTranscoded event", event.VideoID)
}

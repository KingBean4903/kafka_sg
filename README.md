# 🧬 Kafka Schema Registry with Go

A production-grade example project demonstrating how to use **Apache Kafka**, **Confluent Schema Registry**, and **Avro** with **Go**, using the `segmentio/kafka-go` and `hamba/avro` libraries.

---

## 📦 Features

- ✅ Kafka producer written in Go using Avro serialization
- ✅ Schema registration with Confluent Schema Registry
- ✅ Subject name strategy: `TopicNameStrategy`
- ✅ Docker Compose setup for local development
- ✅ Minimal Go consumer to demonstrate decoding with Schema ID

---

## 🧱 Tech Stack

| Component        | Version / Library                     |
|------------------|----------------------------------------|
| Kafka            | bitnami/kafka                          |
| Schema Registry  | Confluent Schema Registry              |
| Go               | 1.24+                                  |
| Kafka client     | [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go) |
| Avro             | [`hamba/avro`](https://github.com/hamba/avro) |
| srclient         | [`riferrei/srclient`](https://github.com/riferrei/srclient) |

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/KingBean4903/kafka_sg.git
cd kafka_sg

```
### 2. Start Kafka + Schema Registry

```bash

docker-compose up -d

```

### Project Structure
.
├── consumer/
│   ├── main.go           
│   └── Dockerfile            
├── docker-compose.yml
├── producer/
│   └── main.go
└── README.md

### References

- Confluent Schema Registry Docs

- Avro Specification

- Hamba Avro Go Library

- Kafka Go Client
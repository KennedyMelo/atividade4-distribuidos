# mini_rabbitmq_benchmark

## Overview
This project provides a minimal implementation of a FIFO queue server over plain TCP, adhering to the PUSH/PULL protocol. It includes producer and consumer clients for benchmarking purposes. Additionally, it offers an alternative implementation using RabbitMQ (via the `streadway/amqp` library) to compare performance metrics.

The repository is designed to help evaluate the latency and throughput of different queueing mechanisms. Each consumer includes a benchmarking helper to capture latency statistics directly from the terminal.

## Directory Structure
```
.
├── simplified
│   ├── server.go       # Minimal in-memory FIFO queue server
│   ├── producer.go     # Producer client for the simplified server
│   └── consumer.go     # Consumer client for the simplified server
└── rabbitmq
    ├── producer_rmq.go # Producer client for RabbitMQ
    └── consumer_rmq.go # Consumer client for RabbitMQ
```

## Prerequisites
- Go 1.24.3 or later
- RabbitMQ broker running and reachable at `amqp://guest:guest@localhost:5672/`

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd atividade4-distribuidos
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

### Simplified Server
1. Start the server:
   ```bash
   go run simplified/server.go
   ```
2. Run the consumer to collect latency stats:
   ```bash
   go run simplified/consumer.go -n 10000
   ```
3. Run the producer to publish timestamps:
   ```bash
   go run simplified/producer.go -n 10000
   ```

### RabbitMQ
1. Ensure RabbitMQ is running and accessible.
2. Run the RabbitMQ consumer:
   ```bash
   go run rabbitmq/consumer_rmq.go -n 10000
   ```
3. Run the RabbitMQ producer:
   ```bash
   go run rabbitmq/producer_rmq.go -n 10000
   ```

## Performance Evaluation
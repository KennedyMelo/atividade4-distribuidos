// mini_rabbitmq_benchmark
// =====================================
// This canvas includes a minimal FIFO queue server implemented over plain TCP, matching the
// PUSH/PULL protocol required by the assignment, alongside producer/consumer clients and an
// alternative producer/consumer pair that targets the commercial RabbitMQ broker (streadway/amqp).
// A small benchmarking helper is embedded in each consumer so you can capture latency statistics
// right from the terminal. Feel free to split these files into separate directories when compiling;
// they live together here for readability.

/*
Directory layout (suggested)
.
├── simplified
│   ├── server.go
│   ├── producer.go
│   └── consumer.go
└── rabbitmq
    ├── producer_rmq.go
    └── consumer_rmq.go

Build & run example (all in different terminals):

Terminal 1 – simplified server
$ go run simplified/server.go

Terminal 2 – simplified consumer (collects latency stats)
$ go run simplified/consumer.go -n 10000

Terminal 3 – simplified producer (publishes timestamps)
$ go run simplified/producer.go -n 10000

Repeat the same pattern with the RabbitMQ versions (be sure the RabbitMQ
broker is running and reachable at amqp://guest:guest@localhost:5672/):

$ go run rabbitmq/consumer_rmq.go -n 10000
$ go run rabbitmq/producer_rmq.go -n 10000

Compare the average/median latency printed by each consumer to complete the
performance evaluation.

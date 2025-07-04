package main

import (
    "flag"
    "fmt"
    "log"
    "math"
    "strconv"
    "time"

    "github.com/streadway/amqp"
)

func main() {
    n := flag.Int("n", 1000, "number of messages to expect before exiting")
    flag.Parse()

    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("channel: %v", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare("bench_queue", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("queue: %v", err)
    }

    msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("consume: %v", err)
    }

    var latencies []time.Duration
    for msg := range msgs {
        sent, _ := strconv.ParseInt(string(msg.Body), 10, 64)
        lat := time.Since(time.Unix(0, sent))
        latencies = append(latencies, lat)
        if len(latencies) >= *n {
            break
        }
    }

    var sum time.Duration
    min := time.Duration(math.MaxInt64)
    max := time.Duration(0)
    for _, l := range latencies {
        if l < min {
            min = l
        }
        if l > max {
            max = l
        }
        sum += l
    }
    avg := sum / time.Duration(len(latencies))
    fmt.Printf("received %d msgs\nmin=%v max=%v avg=%v\n", len(latencies), min, max, avg)
}

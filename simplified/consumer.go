package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "math"
    "net"
    "strconv"
    "strings"
    "time"
)

func main() {
    addr := flag.String("addr", "localhost:9000", "queue server <host:port>")
    n := flag.Int("n", 1000, "number of messages to expect before exiting")
    flag.Parse()

    conn, err := net.Dial("tcp", *addr)
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    var latencies []time.Duration
    for len(latencies) < *n {
        // request next message
        fmt.Fprintf(writer, "PULL\n")
        writer.Flush()

        line, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        line = strings.TrimSpace(line)

        switch {
        case strings.HasPrefix(line, "MSG "):
            tsStr := strings.TrimPrefix(line, "MSG ")
            sent, _ := strconv.ParseInt(tsStr, 10, 64)
            lat := time.Since(time.Unix(0, sent))
            latencies = append(latencies, lat)
        case line == "EMPTY":
            time.Sleep(100 * time.Microsecond) // brief back‑off
        default:
            log.Printf("unexpected: %q", line)
        }
    }

    // simple stats
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
    fmt.Printf("received %d messages\nmin=%v  max=%v  avg=%v\n", len(latencies), min, max, avg)
}
### Estratégia geral

1. **Deixe o consumidor imprimir só o valor médio de latência** (ou um JSON fácil de
   capturar).
2. **Automatize 30 execuções** em um script bash (ou PowerShell) que:

   * inicialize o cenário (servidor ou broker);
   * dispare consumidor → produtor;
   * capture a média de cada rodada;
   * some tudo e divida por 30 ao final.

Abaixo mostro a forma rápida com *bash*, seguida de uma variante em Go caso queira
tudo na mesma linguagem.

---

## 1. Ajuste mínimo no consumidor

Altere o `fmt.Printf` final (tanto no `simplified/consumer.go` quanto no
`rabbitmq/consumer_rmq.go`) para devolver **só a média** ― ou um JSON:

```go
// no lugar de:
fmt.Printf("received %d messages\nmin=%v  max=%v  avg=%v\n", len(latencies), min, max, avg)

// coloque por ex.:
fmt.Println(avg.Microseconds()) // imprime só micro-segundos

// …ou JSON:
fmt.Printf("{\"avg_us\": %d}\n", avg.Microseconds())
```

Isso facilita a captura via `awk`, `jq`, etc.

---

## 2. Script bash para 30 execuções

```bash
#!/usr/bin/env bash
set -euo pipefail
N_RUNS=30
MSG_COUNT=10000

# ---------- fila simplificada ----------
echo "### Simplified FIFO"
sum=0
for i in $(seq 1 $N_RUNS); do
  # 1) inicia servidor em background
  go run simplified/server.go &
  srv_pid=$!
  sleep 0.2  # pequeno tempo para subir

  # 2) inicia consumidor em background (capturando saída)
  avg=$(go run simplified/consumer.go -n $MSG_COUNT &)

  # 3) inicia produtor (bloqueia até terminar)
  go run simplified/producer.go -n $MSG_COUNT

  # 4) espera consumidor terminar
  wait %%% 2>/dev/null || true

  # 5) encerra servidor
  kill $srv_pid

  echo "Run $i => $avg µs"
  sum=$((sum + avg))
done
echo "Média global: $((sum / N_RUNS)) µs"

# ---------- RabbitMQ ----------
echo -e "\n### RabbitMQ (broker externo precisa estar rodando)"
sum=0
for i in $(seq 1 $N_RUNS); do
  avg=$(go run rabbitmq/consumer_rmq.go -n $MSG_COUNT &)
  go run rabbitmq/producer_rmq.go -n $MSG_COUNT
  wait %%% 2>/dev/null || true
  echo "Run $i => $avg µs"
  sum=$((sum + avg))
done
echo "Média global: $((sum / N_RUNS)) µs"
```

> *Dicas*
> • garanta permissão de execução: `chmod +x bench.sh`
> • rode: `./bench.sh`

Se escolheu JSON, troque `avg=$(…)` por algo como:

```bash
avg=$(go run … | jq '.avg_us')
```

---

## 3. Variante puramente em Go

Se quiser evitar scripts externos:

```go
// bench/main.go
package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os/exec"
    "strconv"
)

func run(cmd string, args ...string) (int64, error) {
    c := exec.Command(cmd, args...)
    var out bytes.Buffer
    c.Stdout = &out
    if err := c.Run(); err != nil { return 0, err }
    scanner := bufio.NewScanner(&out)
    scanner.Scan()
    return strconv.ParseInt(scanner.Text(), 10, 64)
}

func main() {
    const runs, msgs = 30, 10000
    var sum int64
    for i := 0; i < runs; i++ {
        // start server
        srv := exec.Command("go", "run", "simplified/server.go")
        if err := srv.Start(); err != nil { panic(err) }
        // consumer (async)
        avgChan := make(chan int64)
        go func() {
            avg, _ := run("go", "run", "simplified/consumer.go", "-n", fmt.Sprint(msgs))
            avgChan <- avg
        }()
        // producer (sync)
        _ = exec.Command("go", "run", "simplified/producer.go", "-n", fmt.Sprint(msgs)).Run()
        avg := <-avgChan
        fmt.Printf("run %d => %d µs\n", i+1, avg)
        sum += avg
        _ = srv.Process.Kill()
    }
    fmt.Printf("GLOBAL AVG: %d µs\n", sum/int64(runs))
}
```

Compile: `go run bench/main.go`.

---

### Por que 30 rodadas?

* Diminui influência de *jitters* (GC, SO, CPU turbo, etc.).
* Permite extrair desvio-padrão se quiser um teste estatístico
  (ex.: **t-test** para significância entre as duas médias).

Se precisar de:

* **CSV** para gráficos → acrescente `echo "$i,$avg" >> result.csv`.
* **Desvio padrão / mediana** → armazene em slice e processe no fim.
* **Gráficos** → use Python/Excel ou `gnuplot` a partir do CSV.

Pronto! Assim você roda cada versão 30 vezes, captura a latência média de
cada execução e calcula o “average of averages” automaticamente.

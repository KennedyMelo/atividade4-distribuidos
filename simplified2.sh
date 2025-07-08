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
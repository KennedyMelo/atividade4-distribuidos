#!/usr/bin/env bash
set -euo pipefail
N_RUNS=30
MSG_COUNT=10000

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
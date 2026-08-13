#!/usr/bin/env bash
set -euo pipefail

if command -v ss >/dev/null 2>&1; then
  echo "Все TCP-соединения:"
  ss -tn
  echo
  echo "ESTABLISHED:"
  ss -tn state established
  echo
  echo "TIME-WAIT:"
  ss -tn state time-wait
  echo
  echo "Слушающие TCP-порты:"
  ss -tln
elif command -v netstat >/dev/null 2>&1; then
  echo "TCP-соединения в состояниях LISTEN, ESTABLISHED и TIME_WAIT:"
  netstat -an -p tcp | awk 'NR <= 2 || /LISTEN|ESTABLISHED|TIME_WAIT/'
else
  echo "Не найдено ни ss, ни netstat" >&2
  exit 1
fi

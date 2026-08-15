#!/usr/bin/env bash
set -euo pipefail

port="${PORT:-9090}"

if command -v ss >/dev/null 2>&1; then
	echo "TCP-соединения для порта ${port}:"
	ss -tan | awk -v port=":${port}" 'NR == 1 || index($4, port) || index($5, port)'
elif command -v netstat >/dev/null 2>&1; then
	echo "TCP-соединения для порта ${port}:"
	netstat -an -p tcp | awk -v dot_port=".${port}" -v colon_port=":${port}" \
		'NR <= 2 || index($4, dot_port) || index($5, dot_port) || index($4, colon_port) || index($5, colon_port)'
else
	echo "Не найдено ни ss, ни netstat" >&2
	exit 1
fi

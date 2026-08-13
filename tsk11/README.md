# Практика: стандартная сетевая библиотека Go

Все программы запускаются независимо из корня репозитория.

## URL Shortener

```bash
go run ./project_11/cmd/url_shortener
curl -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://golang.org"}'
curl http://localhost:8080/stats
```

Для проверки перенаправления подставьте полученный код:

```bash
curl -v http://localhost:8080/s/ABC123
```

## HTTP-клиент с retry

```bash
go run ./project_11/cmd/retry_client
go run ./project_11/cmd/retry_client -url https://httpbin.org/status/500 -retries 2
```

## Проверка доступности сервисов

Без аргументов проверяется встроенный список URL. Можно передать собственный список:

```bash
go run ./project_11/cmd/service_checker
go run ./project_11/cmd/service_checker https://go.dev https://github.com
```

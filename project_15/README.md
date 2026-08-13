# Практика: HTTP и HTTPS

Все команды выполняются из корня репозитория.

## Исследование HTTP с curl

```bash
./project_16/http_requests.sh
```

По умолчанию используется `https://httpbin.org`. Другую совместимую базу можно передать через переменную `HTTPBIN_BASE`.

## HTTP-клиент

```bash
go run ./project_16/cmd/http_client
go run ./project_16/cmd/http_client -url https://example.com -timeout 5s
```

## Файловый сервер с CORS

```bash
go run ./project_16/cmd/file_server
curl -i http://localhost:8080/
curl -i -X OPTIONS http://localhost:8080/
```

Для `OPTIONS` сервер возвращает статус `204`, а для файлов — заголовок `Access-Control-Allow-Origin: *`.

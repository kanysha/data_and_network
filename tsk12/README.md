# Практика: HTTP и HTTPS

Все команды выполняются из корня репозитория.

## Исследование HTTP с curl

```bash
./tsk12/http_requests.sh
```

По умолчанию используется `https://httpbin.org`. Другую совместимую базу можно передать через переменную `HTTPBIN_BASE`.
Таймаут каждого запроса по умолчанию равен 15 секундам и настраивается через `HTTP_TIMEOUT`.

## HTTP-клиент

```bash
go run ./tsk12/cmd/http_client
go run ./tsk12/cmd/http_client -url https://example.com -timeout 5s
```

## Файловый сервер с CORS

```bash
go run ./tsk12/cmd/file_server
curl -i http://localhost:8080/
curl -i -X OPTIONS http://localhost:8080/
```

Для `OPTIONS` сервер возвращает статус `204`, а для файлов — заголовок `Access-Control-Allow-Origin: *`.
Пример фактических HTTP-ответов сохранён в `results.txt`.

# ПЗ №7 — Redis, SET/GET и TTL

Сортов Семён Александрович, ЭФМО-01-26. Дата: 25.10.2026.

Простой HTTP-доступ к кэшу Redis: /set сохраняет строку с TTL 10 секунд, /get читает значение, /ttl показывает оставшееся время. Код соответствует структуре из задания: cmd/server и internal/cache.

Требования: Go 1.26+, Docker Desktop, PowerShell или curl. Проверено с Go 1.27.2 на Windows 11 amd64 и Redis 7.

## Запуск

```powershell
docker compose up -d
docker compose exec redis redis-cli PING
go run ./cmd/server
```

Redis доступен на localhost:6379, HTTP — на localhost:8080. Практики с одним HTTP-портом запускаются отдельно.

## Проверка в другом окне

```powershell
curl.exe 'http://localhost:8080/set?key=test&value=hello'
# OK: test=hello (TTL 10s)
curl.exe 'http://localhost:8080/get?key=test'
# VALUE: test=hello
curl.exe 'http://localhost:8080/ttl?key=test'
# TTL for test: 10s (или меньше, в зависимости от времени запроса)
docker compose exec redis redis-cli TTL test
Start-Sleep -Seconds 11
curl.exe -i 'http://localhost:8080/get?key=test'
# HTTP 404: key not found
docker compose exec redis redis-cli EXISTS test
# 0
```

Если ключ отсутствует, /ttl показывает key not found (-2), что соответствует ответу команды TTL Redis. Пустой key или value в /set даёт 400; /get и /ttl также требуют key. Ошибка связи с Redis возвращает 500, отсутствие значения — 404.

## Проверки и типовые проблемы

```powershell
go test ./...
go vet ./...
go build -o bin/app.exe ./cmd/server
```

- Нет связи с Redis: проверить Docker Desktop и ответ PING.
- 404 после паузы: ключ штатно удалён по TTL; заполнить кэш снова.
- Порт 8080 занят: остановить ранее запущенный сервер другой практики.

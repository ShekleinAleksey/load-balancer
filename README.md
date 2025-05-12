# Load Balancer with Rate Limiting

Go-реализация балансировщика нагрузки с rate limitting и health checks.

## Features
- Алгоритмы балансировки: Round-Robin, Least Connections, Random
- Rate limiting (Token Bucket)
- Health checks бэкендов

## Быстрый старт

### Запуск с Docker
```bash
docker-compose up --build
```

### Тестирование
```bash
# Unit-тесты
go test -v ./...


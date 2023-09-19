# Тестовое задание от компании Effective Mobile

## Описание

С заданием можно ознакомиться в файле task.pdf.

В примере используются Kafka, Postgres, Redis, Rest API, GraphQL API и получение данных со сторонних сервисов для обогащения данных.

Для генерации тестовых записей можно запустить приложение producer, которое будет генерить пользователей со случайными именами и отправлять их в топик кафки.

## Запуск

### Запуск Kafka, Postgres, Redis

```bash
docker-compose up -d
```

### Запуск consumer для получения сообщений из Кафки

```bash
go run ./cmd/consumer/.
```

### Запуск producer для генерации сообщений в Кафку

```bash
go run ./cmd/producer/.
```

### Запуск Rest API

```bash
go run ./cmd/rester/.
```

### Запуск GraphQL API

```bash
go run ./cmd/grapher/.
```

## Остановка Kafka, Postgres, Redis

```bash
docker-compose down
```

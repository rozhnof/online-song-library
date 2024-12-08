# Song Service

## Requirements

- Docker
- Docker Compose
- Make

## Build and Run with Docker

Для запуска сервиса с использованием Docker выполните следующие шаги:

1. Установите необходимые переменные среды:
    ```bash
    make env
    ```
2. Соберите Docker образ:
    ```bash
    make build-service
    ```
3. Поднимите сервис:
    ```bash
    make up-service
    ```

## Stopping the Service

Для остановки сервиса выполните команду:
```bash
make down-service
```

## API Documentation

Интерактивная документация с использованием Swagger UI доступна по следующему адресу:

[![Swagger UI](https://img.shields.io/badge/-API%20Documentation-blue)](http://localhost:8080/swagger/index.html)

## Running Tests

### Functional Tests

Для запуска функциональных тестов сначала настройте тестовое окружение:

1. Установите необходимые тестовые переменные среды:
    ```bash
    make test-env
    ```
2. Соберите Docker образ для тестового окружения:
    ```bash
    make build-test-service
    ```
3. Поднимите сервис для тестов:
    ```bash
    make up-test-service
    ```

Запуск функциональных тестов:

```bash
make run-functional-tests
```

Для остановки тестового сервиса используйте команду:

```bash
make down-test-service
```

## Database Migrations

### Running Migrations

Для применения миграций выполните команду:

```bash
make migration-up
```

### Rolling Back Migrations

Для отката миграций выполните команду:

```bash
make migration-down
```

### Creating a New Migration

Для создания новой миграции выполните команду:

```bash
make migration-create name=<migration_name>
```
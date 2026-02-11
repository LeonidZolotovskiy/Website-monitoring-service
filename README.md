# Website-monitoring-service
# Site Monitor

Простой мониторинг сайтов на Go.

## Описание

Эта программа проверяет работу заданного списка сайтов каждую минуту.  
Проверка осуществляется по HTTP-статусу: если сайт отвечает `200 OK` — считается работающим, иначе — не работает.

Программа выводит в консоль:

Site https://www.google.com ok
Site https://www.nonexistentsite12345.com NOT ok

## Использование

Запустите программу:

go run ./cmd/monitor -config ./configs/sites.yaml

Настройка списка сайтов

Список сайтов хранится в массиве sites в файле main.go.
Пример:

sites := []string{
    "https://www.google.com",
    "https://www.github.com",
    "https://www.python.org",
    ...
}

Особенности

Таймаут запроса 5 секунд.

Проверка всех сайтов повторяется каждые 60 секунд.

Можно легко расширить список сайтов.

### Run migrations

Export database url:

export DB_URL="postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

Apply migrations:

migrate -path migrations -database "$DB_URL" up

Rollback last migration:

migrate -path migrations -database "$DB_URL" down 1

#### Запуск тестов

1. Запуск всех тестов
make test

Или напрямую через Go:

go test ./...

2. Подробный вывод тестов
make test-verbose

Или:

go test -v ./...

3. Проверка покрытия кода
make test-cover

Выведет процент покрытия и создаст файл coverage.out.

4. Генерация HTML-отчета по покрытию
make test-cover-html

Создаётся файл coverage.html.

Его можно открыть в браузере для визуального анализа покрытия:

open coverage.html  # macOS
xdg-open coverage.html  # Linux

Примеры флагов go test

-run TestName — запуск конкретного теста:

go test -v ./... -run TestCheckSite

-short — пропуск долгих тестов:

go test -v -short ./...

-coverprofile и -cover — проверка покрытия:

go test -cover -coverprofile=coverage.out ./...
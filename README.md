# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите Go версии 1.25.1.
2. Склонируйте репозиторий.
3. Подтяните зависимости:

```bash
go mod tidy
```

4. Создайте в корне `.env.override`:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

5. Запустите тесты:

```bash
go test -v ./...
```

## Запуск линтера

```bash
golangci-lint run
```

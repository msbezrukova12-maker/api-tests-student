# API-тесты TestBoard — лабораторная Avito

Тесты работают с реальным [TestBoard](https://testboard.avito.com/).
[Swagger](https://testboard.avito.com/swagger/index.html) — источник контрактов API.

## Подготовка

Нужны Git, Go 1.25.1 или новее и golangci-lint v2, собранный с поддержкой установленной версии Go.
[Установка Go](https://go.dev/doc/install), [установка линтера](https://golangci-lint.run/docs/welcome/install/local/).

```bash
git clone https://github.com/msbezrukova12-maker/api-tests-student.git
cd api-tests-student
git switch main
go mod download
```

Создайте в корне `.env.override` с данными зарегистрированного пользователя TestBoard:

```dotenv
TEST_LOGIN=your_email
TEST_PASSWORD=your_password
```

Базовый адрес `API_URL=https://testboard.avito.com/api/v1` уже указан в `.env`.
`.env.override` загружается после `.env` и исключён из Git. Не публикуйте логин, пароль и токены.
Нужен сетевой доступ к TestBoard; для favorites тесты сами регистрируют отдельного пользователя через API.

## Запуск

Из корня репозитория:

```bash
go test -p 1 -count=3 ./...
golangci-lint run
```

`-p 1` запускает пакеты последовательно, `-count=3` проверяет три повторения без кэша.
Для подробного вывода шагов добавьте `-v`. Запуск отдельного сценария:

```bash
go test -p 1 -count=1 -v ./tests/scenarios/advertisementCreate -run '^TestCreateAdvertisementAllFields$'
```

Тесты выводят `Precondition`, `Step`, `Result`. Для отладки HTTP-запросов можно добавить
`DEBUG=true` в `.env.override`. Такой вывод может содержать авторизационные данные, не публикуйте его.

## Покрытие и структура

- Создание объявления со всеми полями и тремя фото, сравнение параметров через GET, поиск по полному title, получение фото.
- Невалидные поля создания: обязательные поля, длина title/description, цена, quantity, отсутствие фото; точные `400`, `error`, `message`.
- Отсутствующий и неверный токен: точные `401` и сообщения для создания и всех операций favorites.
- Favorites: добавление, просмотр, удаление из списка, исчезновение после удаления самого объявления.
- Несуществующее объявление: точные `404` и сообщения добавления/удаления.
- POST favorites: `400` без обязательного ID, `403` для своего объявления, `409` при повторном добавлении.
- Получение собственных объявлений с данными, созданными самим тестом, и проверка неверного токена.

Архитектура: HTTP-клиенты в `internal/client/http` → проверки и модели в `internal/managers` → сценарии в `tests/scenarios`.
Общие данные и подготовка находятся в `tests/fixtures`, проверки полей — в `internal/assertions`.

Фактический Swagger использует `GET /advertisement?id={id}` и `DELETE /advertisement?id={id}`;
фото запрашиваются через `GET /advertisements/{id}/photos`. Это отличается от записи пути в тексте лабораторной.
Статусы и сообщения сверены с TestBoard. Ошибки `5xx` приводят к падению, повторные попытки и пропуски их не скрывают.

Каждый сценарий создаёт свои объявления с уникальными названиями и удаляет их через `t.Cleanup`.
Favorites используют отдельного нового пользователя с пустым избранным, поэтому не затрагивают избранное основной учётной записи.
API не предоставляет удаления пользователей: зарегистрированные тестовые учётные записи остаются на стенде.

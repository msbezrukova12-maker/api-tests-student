# API-тесты для проекта https://testboard.avito.com/

Проект содержит API-тесты TestBoard, построенные по архитектуре **Client Layer → Manager Layer → Scenario**.

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Библиотека apitest](https://github.com/steinfletcher/apitest)
3. [Библиотека gjson](https://github.com/tidwall/gjson)

## Требования

- Go 1.25.1 или выше
- Git
- golangci-lint
- учётная запись TestBoard

## Настройка

1. Склонируйте fork:

```bash
git clone https://github.com/<ваш-логин>/api-tests-student.git
cd api-tests-student
```

2. Установите зависимости:

```bash
go mod tidy
```

3. Создайте в корне файл `.env.override`:

```env
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

Файл `.env.override` находится в `.gitignore` и не должен попадать в Git.

Базовый URL уже задаётся в `.env`:

```env
API_URL=https://testboard.avito.com/api/v1
```

## Запуск тестов

Все тесты:

```bash
go test -v ./...
```

Проверка идемпотентности тремя последовательными прогонами:

```bash
go test -count=3 ./...
```

Только создание объявления:

```bash
go test -v ./tests/scenarios/advertisementCreate/
```

Только избранное:

```bash
go test -v ./tests/scenarios/favorites/
```

Поиск объявлений:

```bash
go test -v ./tests/scenarios/advertisementSearch/
```

## Линтер

Установите golangci-lint по официальной инструкции и запустите из корня проекта:

```bash
golangci-lint run
```

## Что покрыто

### Создание объявления

- создание объявления со всеми параметрами;
- получение созданного объявления по ID;
- поиск по полному title;
- проверка фотографий через `GET /advertisements/{id}/photos`;
- обязательные поля и ошибки валидации;
- пустой и слишком длинный title;
- отрицательная цена;
- отсутствие фотографий;
- отсутствие и невалидный bearer-токен.

### Избранное

- добавление объявления в избранное и проверка через `GET /favorites`;
- удаление объявления из избранного;
- каскадное исчезновение favorite после удаления самого объявления;
- добавление и удаление несуществующего `ad_id`;
- отсутствие и невалидный bearer-токен;
- повторное добавление одного объявления;
- запрет добавления собственного объявления.

## Изоляция тестов

Тесты сами создают необходимые объявления и регистрируют cleanup через `t.Cleanup`.
Для сценариев избранного автоматически создаётся отдельный тестовый пользователь, поэтому ревьюеру не нужен второй аккаунт.
Перед каждым favorites-сценарием его список избранного очищается, а после сценария выполняется повторный cleanup.

Тесты не зависят от заранее созданных объявлений в аккаунте и рассчитаны на повторные запуски.

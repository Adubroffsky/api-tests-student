# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Библиотека для анализа JSON в Go](https://github.com/tidwall/gjson)

## Требования

Для запуска проекта необходимы:

- Go версии 1.25.1 или выше;
- Git;
- IDE для разработки, например [GoLand](https://www.jetbrains.com/go/) или [VS Code](https://code.visualstudio.com/).

## Подготовка проекта

1. Склонируйте репозиторий:

```bash
git clone https://github.com/<your-username>/api-tests-student.git
```

2. Перейдите в каталог проекта:

```bash
cd api-tests-student
```

3. Установите зависимости:

```bash
go mod tidy
```

## Настройка переменных окружения

В корне проекта создайте файл `.env.override`.

Укажите в нём данные двух тестовых пользователей:

```env
TEST_LOGIN=your_login
TEST_PASSWORD=your_password

SECOND_TEST_LOGIN=your_second_login
SECOND_TEST_PASSWORD=your_second_password
```

- `TEST_LOGIN` и `TEST_PASSWORD` — логин и пароль первого тестового пользователя.
- `SECOND_TEST_LOGIN` и `SECOND_TEST_PASSWORD` — логин и пароль второго тестового пользователя.

Для локального запуска можно также использовать файл `.env`. При загрузке настроек сначала читается `.env`, затем `.env.override`.

**Важно:** не публикуйте логины и пароли и не добавляйте `.env.override` в Git.

## Запуск тестов

Для запуска всех тестов из корня проекта:

```bash
go test -p 1 -count=1 -v ./...
```

Параметр `-p 1` запускает пакеты последовательно, а `-count=1` отключает использование кэша результатов тестирования.

Для запуска отдельных наборов тестов:

**Объявления:**

```bash
go test -v ./tests/scenarios/advertisement/
```

**Избранное:**

```bash
go test -v ./tests/scenarios/favorites/
```

**Мои объявления:**

```bash
go test -v ./tests/scenarios/myAdvertisement/
```

## Тестовая фотография

Файл `test_photo.webp` используется тестами при создании объявлений.

Он должен находиться в корне проекта. Не перемещайте и не удаляйте его.

## Логирование

Для включения логирования HTTP-запросов и ответов в консоль добавьте в `.env.override`:

```env
DEBUG=true
```

## Запуск линтера

1. Установите [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/).
2. Запустите проверку из корня проекта:

```bash
golangci-lint run
```
# Backend Hackathon - Черная Жемчужина

Техническое описание Go-сервера для геймифицированного образовательного сервиса.

## 🚀 О проекте

REST API для мобильного приложения с геймификацией обучения финансовой грамотности. Пользователи выращивают виртуального питомца, проходя обучающие материалы и тесты.

## 🛠️ Технологии

- **Язык:** Go 1.21\+
- **Фреймворк:** Gin
- **База данных:** PostgreSQL
- **Контейнеризация:** Docker \+ Docker-compose
- **Документация API:** Swagger 2.0 
- **Миграции:** Нативные SQL

## 📚 Swagger Документация

Автоматически сгенерированная документация доступна по адресу:

<http://localhost:8080/swagger/index.html>

### Основные DTO модели

#### Pet (Питомец)

```
{
  "id": 1,
  "name": "Медвежонок",
  "age": 2,
  "exp": 150
}
```

#### **Prize (Призы)**

```
{
  "type": "кешбек",
  "title": "Кешбек 5%",
  "description": "Кешбек на следующие покупки"
}
```

#### **Quiz (Тесты)**

```
{
  "id": 1,
  "title": "Финансовый тест",
  "content": "Как правильно планировать бюджет?",
  "options": ["Вариант 1", "Вариант 2", "Вариант 3"],
  "correctAnswer": "1"
}
```

## **🔌 API Endpoints**

### **🐻 Pet Management**

#### **`POST /pet/name`**

**Переименовать питомца**

**Request Body:**

```
{
  "userID": 123,
  "name": "Медвежонок"
}
```

**Responses:**

* `200 OK` - Успешное выполнение

* `400 Bad Request` - Неверный запрос

* `500 Internal Server Error` - Ошибка сервера

#### **`POST /pet/xp`**

**Добавить опыт питомцу**

**Request Body:**

```
{
  "userID": 123,
  "exp": 100
}
```

**Responses:**

* `200 OK` - Успешное выполнение

* `400 Bad Request` - Неверный запрос

* `500 Internal Server Error` - Ошибка сервера

#### **`GET /pet/{id}`**

**Получить питомца по ID пользователя**

**Parameters:**

* `id` (integer, path) - User ID

**Response:**

```
{
  "id": 1,
  "name": "Медвежонок",
  "age": 2,
  "exp": 150
}
```

**Responses:**

* `200 OK` - Успешное выполнение

* `400 Bad Request` - Неверный запрос

* `404 Not Found` - Питомец не найден

* `500 Internal Server Error` - Ошибка сервера

### **🎁 Prizes Management**

#### **`POST /prizes/{id}/available`**

**Получить доступные призы для пользователя**

**Parameters:**

* `id` (integer, path) - User ID

**Response:**

```
{
  "prizes": [
    {
      "type": "кешбек",
      "title": "Кешбек 5%",
      "description": "Кешбек на следующие покупки"
    }
  ]
}
```

#### **`GET /prizes/{id}/my`**

**Получить призы пользователя**

**Parameters:**

* `id` (integer, path) - User ID

**Response:** Аналогично `/prizes/{id}/available`

### **📚 Learning Content**

#### **`GET /sections`**

**Получить все секции с айтемами**

**Response:**

```
[
  {
    "id": 1,
    "title": "Основы финансов",
    "items": [
      {
        "itemID": 1,
        "sectionID": 1,
        "title": "Что такое бюджет",
        "isTest": false
      }
    ]
  }
]
```

#### **`POST /sections`**

**Создать новую секцию**

**Request Body:**

```
{
  "title": "Новая секция"
}
```

**Response:** `201 Created` с данными секции

#### **`GET /sections/{id}/items`**

**Получить айтемы конкретной секции**

**Parameters:**

* `id` (integer, path) - Section ID

**Response:** Массив элементов секции

#### **`POST /sections/{id}/items`**

**Создать айтем в секции**

**Parameters:**

* `id` (integer, path) - Section ID

**Request Body:**

```
{
  "sectionId": 1,
  "itemId": 1,
  "title": "Новый элемент",
  "isTest": true
}
```

**Response:** `201 Created` с данными элемента

### **🎯 Quiz Management**

#### **`GET /quiz/{id}`**

**Получить квиз по ID**

**Parameters:**

* `id` (integer, path) - Quiz ID

**Response:**

```
{
  "id": 1,
  "title": "Финансовый тест",
  "content": "Как правильно планировать бюджет?",
  "options": ["Вариант 1", "Вариант 2", "Вариант 3"],
  "correctAnswer": "1"
}
```

### **📖 Theory Management**

#### **`GET /theory/{id}`**

**Получить теорию по ID**

**Parameters:**

* `id` (integer, path) - Theory ID

**Response:**

```
{
  "id": 1,
  "title": "Основы бюджетирования",
  "content": "Бюджет - это план доходов и расходов..."
}
```

#### **`POST /theory`**

**Создать теорию**

**Request Body:**

```
{
  "id": 1,
  "title": "Новая теория",
  "content": "Содержание теории..."
}
```

**Response:** `201 Created` с данными теории

### **✅ Daily Tasks**

#### **`GET /tasks/daily`**

**Получить ежедневные задания**

**Response:**

```
[
  {
    "ID": 1,
    "Title": "Пройдите первый тест"
  }
]
```

## **🗄️ База данных**

### **Миграции**

* `migrations/` - SQL файлы миграций

* Автоматическое применение при старте

* Версионное управление схемой

### **Основные таблицы**

* `users` - Пользователи

* `pets` - Питомцы

* `prizes` - Призы

* `sections` - Разделы обучения

* `section_items` - Элементы разделов

* `theory` - Теоретические материалы

* `quiz` - Тесты и викторины

* `tasks` - Ежедневные задания

## **🚀 Запуск проекта**

### **Требования**

* Docker & Docker-compose

* Go 1.21\+

### **Локальная разработка**

1. **Клонирование репозитория:**

```
git clone <repository-url>
cd backend-hackathon
```

2. **Настройка окружения:**

```
cp .env.example .env
# Отредактируйте .env файл
```

3. **Запуск контейнеров:**

```
docker-compose up -d
```

4. **Запуск приложения:**

```
go run cmd/server/main.go
```

### **Доступ к сервисам**

* **API:** [http://localhost:8080](http://localhost:8080/)

* **Swagger UI:** <http://localhost:8080/swagger/index.html>

* **PostgreSQL:** localhost:5432

## **🔧 Конфигурация**

Основные параметры в `.env` файле:

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=hackathon
DB_USER=user
DB_PASSWORD=password
SERVER_PORT=8080
```

## **🐛 Обработка ошибок**

Все ошибки возвращаются в формате:

```
{
  "error": "Описание ошибки"
}
```

**Статус коды:**

* `200` - Успешный запрос

* `201` - Успешное создание

* `400` - Неверный запрос

* `404` - Ресурс не найден

* `500` - Внутренняя ошибка сервера

## **📝 Логирование**

Структурированное логирование через кастомный логгер:

* Уровни логирования (DEBUG, INFO, WARN, ERROR)

* JSON формат для production

* Контекстные поля

## **🤝 Разработка**

### **Генерация Swagger документации**

```
# Установка swag
go install github.com/swaggo/swag/cmd/swag@latest

# Генерация документации
swag init -g cmd/server/main.go

# Обновление спецификации
swag fmt
swag init
```

### **Code Style**

* Go fmt

* Go vet

* Стандарты Go Code Review

### **Тестирование**

```
go test ./...
```

### **Миграции**

Добавление новых миграций в папку `migrations/` с префиксами:

* `{timestamp}_description.up.sql` - Применение

* `{timestamp}_description.down.sql` - Откат

## **📄 Лицензия**

Проект распространяется под лицензией MIT. Подробнее в файле LICENSE.

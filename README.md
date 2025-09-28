# Order manager
Микросервис для управления заказами с использованием Go, PostgreSQL и Kafka.
## Функциональность
- Прием сообщений о заказах из Kafka
- Сохранение данных в PostgreSQL
- In-memory кэширование для быстрого доступа
- Восстановление кэша при перезапуске
- HTTP-сервер для получения информации о заказах
- Веб-интерфейс для поиска заказов

## Технологии
- Go 1.23
- PostgreSQL - база данных
- Kafka - брокер сообщений
- Docker - контейнеризация
- Chi - HTTP роутер
- Goose - миграции
- swagger - документация
- testify - тестирование
- go-playground/validator - валидация
- log/slog - логирование
## Структура
```
order-manager/
├── cmd/
│   └── api/
│       └── main.go                 # Точка входа приложения
├── docs/                           # Документация
├── internal/
│   ├── api/                        # Запуск и остановка приложения
│   │   └── api.go                       
│   ├── cache/                      # In-memory кэш
│   │   └── cache.go
│   ├── config/                     # Конфигурация
│   │   └── config.go
│   ├── controller/                 # Слой controller              
│   │   ├── http
|   |   |   ├── handlers.go         # Handlers
|   |   |   └── router.go           # HTTP сервер
│   │   └── kafka
|   |   |   └── consumer.go         # Kafka консьюмер 
│   ├── models/                     # Модели данных
│   │   └── model.go
│   ├── repository/                 # Слой repository
│   │   └── repository.go
│   └── service/                    # Слой service
│       ├── service.go              
│       └── service_test.go
├── migrations/                     # Миграции
├── mocks/                          # Моки
├── prg/
│   ├── db/                         # Соединение с PostgreSQL                                            
│   │   └── db.go                       
│   ├── errorx/                     # Кастомные ошибки
│   │   └── errorx.go               
│   ├── web/                        # Фронтенд
│       └── index.html           
├── producer/                       # Kafka продюсер
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── go.mod                          
├── go.sum                          
├── docker-compose.yml              
├── Dockerfile
└── .env.example                    
```
## Установка и запуск
### 1. Клонирование репозитория  
```bash
git clone https://github.com/EBichuk/order-manager.git
cd order-manager
```
### 2. Запуск приложения
```bash
docker compose up -d
```
Эта команда запустит:
- PostgreSQL на порту 5432
- Сервис на порту 8081
- Kafka с тремя контролерами и тремя брокерами + Kafka ui на порту 8082
- Сервис producer, который отправляет в kafka 20 заказов

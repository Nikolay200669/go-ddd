## go-ddd

Use DDD (Domain-Driven Design) in Go development.


### Project structure
```shell
project
├── domain
│   ├── model.go        // Определение моделей доменных объектов
│   └── repository.go   // Интерфейс репозитория для работы с базой данных
├── infrastructure
│   └── persistence.go  // Реализация репозитория с использованием GORM
├── interfaces
│   └── handlers.go     // Обработчики HTTP запросов с использованием Gin
└── main.go             // Основной файл приложения

```
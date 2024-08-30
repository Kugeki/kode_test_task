# Тестовое задание в KODE

# Проверяющим
Запуск сервера:
```
make run
```

Остановка сервера:
```
make stop
```

При запуске сервера будет доступна страница Swagger на `localhost:8080/`.
Там можно протестировать API (с помощью кнопки `Try it out`).

Также есть автотесты. Они находятся в папке `tests/`. 
Их запуск:
```
make run-tests
```
Все возможности API можно посмотреть в файле `tests/helpers_tests.go` и на странице Swagger.

Пользователей три штуки:
'''
test1:password1
test2:password2
test3:password3
'''

# Задача
Необходимо спроектировать и реализовать на Golang сервис, 
предоставляющий REST API интерфейс с методами:
- добавление заметки
- вывод списка заметок

При сохранении заметок необходимо орфографические ошибки
валидировать при помощи сервиса [Яндекс.Спеллер](https://yandex.ru/dev/speller/) 
(добавить интеграцию с сервисом). Также необходимо реализовать 
аутентификацию и авторизацию. Пользователи должны иметь доступ только к своим
заметкам. Возможность регистрации не обязательна, допустимо иметь
предустановленный набор пользователей (механизм хранения учетных
записей любой, вплоть до hardcode в приложении).

# Условия
Для реализации сервиса использовать язык программирования Go
- Сервис должен работать через REST API, для передачи данных
использовать формат json
- Логирование событий в едином формате
- В целом рекомендуется использовать преимущественно стандартную
библиотеку и библиотеки golang.org , библиотека логгера - по выбору,
библиотека web-сервера - chi, gorilla или стандартная
- Не используем gin, gorm и другие подобные фреймворки или ORM
- Запуск сервиса и требуемой им инфраструктуры должен
производиться в докер контейнерах, необходимо продумать удобство
проверки работоспособности методов API при ревью задачи (шаблоны
curl запросов, postman коллекция, автотесты и т.п.)


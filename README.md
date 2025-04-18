# Часть сервиса аутентификации
Выполнено по данному тз: https://medods.yonote.ru/share/a74f6d8d-1489-4b54-bd82-81af5bf50a03/doc/test-task-backdev-sCBrYs5n6e
##
Конфигурация находится в ```./config/prod.yaml```
```yaml
env: "prod" #Окружение
database: #Параметры для подключения к базе данных
  host: 0.0.0.0 #ip адрес сервера с БД
  port: 5432 #порт для подключения к БД
  user: postgres #пользователь БД
  password: postgres #пароль БД
  name: authdb #база данных
paths: #пути для логов, желательно не трогать
  files: ./files #папка с файлами
  logDir: ./log #каталог для логов
  logName: auth_service.log #имя файлов с логами
jwtSecret: "secret" #секретный ключ для кодирования jwt токена
server: #серверные настройки
  port: 8080 #порт прослушивание
  timeout: 5s #тайм-аут ожидания
```

## Запуск
Присутствует запуск с аргументами: ```service --config=./config/prod.yaml```
## Make
В корне проекта присутствует **makefile**, содержит команды:
```bash
make build  #сборка проекта, бинарник создается в корне директории
make run_local #запуск сервиса с локальной конфигурацией (--config=config/local.yaml)
make run_prod #запуск сервиса с продакшен конфигурацией (--config=config/prod.yaml)
make clear #отчистка локальных папок
```

## Тесты
Тесты находятся в папке: ```tests```

## Docker
В корне присутствует **Dockerfile** и **docker-compose.yml**. По умолчанию будет скопирована конфигурация
из папки ```./config/prod.yaml``` и скопирована в контейнер как ```prod.yaml```, сервис запустится с аргументом
```--config=prod.yaml```

> [!IMPORTANT]
> Если не указать аргументы запуска, по умолчанию будет загружен конфиг по пути ```./config/local.yaml```, если
> конфигурация по данному пути не будет найдена, то сервис попытается загрузить конфигурацию по пути ```./local.
> yaml```, в случае если конфигурация так и не будет найдена, сервис завершится с ошибкой.

> [!IMPORTANT]
> Для использования моего **makefile** обязательно необходима утилита **[make](https://www.make.com/en)**

# Документация
Эндпоинты:
**GET**: ```/api/v1//auth/token?id=```
Возвращает JSON:
```json
    {
        "AccessToken": "",
        "RefreshToken": ""
    }
```

Эндпоинты:
**POST**: ```/api/v1//auth/refresh```
Принимает JSON:
```json
    {
        "AccessToken": "",
        "RefreshToken": ""
    }
```

Возвращает JSON:
```json
    {
        "AccessToken": "",
        "RefreshToken": ""
    }
```
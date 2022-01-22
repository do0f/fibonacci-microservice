# fibonacci-microservice
REST API and gRPC microservice for getting fibonacci numbers sequence

Микросервис возвращает числа фибоначчи с номерами от first до last, нумерация начинается с 0.

Вычисленные числа кэшируются в Redis.

Данные для подключения к сервису задаются в .env файле (Порты серверов, пароль от БД)
Развертка сервиса выполняется с помощью docker-compose

Команды:
- git clone https://github.com/do0f/fibonacci-microservice
- docker-compose up --build fibonacci

Пример работы сервиса:
![Пример](https://github.com/do0f/fibonacci-microservice/blob/main/example.png)

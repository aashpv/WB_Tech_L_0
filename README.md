<div align="center">
  <h1 align="center">WB Order Service</h1>
  <h3>Cервис для управления заказами.</h3>
</div>

<br/>

WB Order Service — это сервис для управления заказами, который включает в себя обработку и хранение информации о заказах, доставке, оплате и товарах. Сервис использует PostgreSQL для хранения данных и NATS Streaming для асинхронного обмена сообщениями.

## Демо

![WB_Order_Service.jpg](https://sun9-47.userapi.com/impg/5Vm9I-W8QYriJLjPFC0o0IwzFBS_BVvEEYoVFw/xgvUC1PRiy8.jpg?size=1510x890&quality=96&sign=905cc8f6536d3bffc68446c28b6fefbd&type=album)

## Стек

- [Go](https://go.dev/) – Programming language
- PostgreSQL - Database
- [NATS Streaming](https://github.com/nats-io/nats-streaming-server) – Platform for asynchronous messaging
- [Go-cache](https://github.com/patrickmn/go-cache) – Key/value storage in memory
- [Chi](https://github.com/go-chi/chi) – Framework

## Приступая к работе

### Предварительные требования

Вот что вам нужно для запуска:

- Go (version >= 15)
- PostgreSQL Database
- NATS Streaming Server

### 1. Склонируйте репозиторий

```shell
git clone https://github.com/aashpv/WB_Tech_L_0
```

### 2. Создайте таблицы

```shell
cd WB_Tech_L_0\scripts
psql -h localhost -p 5432 -U your_username -f create.sql
```

### 3. Настройте config.yaml

```yaml
storage_conn: "user=your_username password=your_password dbname=wb-order-service host=localhost port=5432 sslmode=disable"
http_server:
  address: "localhost:8090"
  timeout: 4s
  idle_timeout: 60s
nats_streaming:
  cluster_id: "test-cluster"
  client_id: "wb-order-service"
  subject: "wb-tech"
```

### 4. Запустите nats-streaming-server

```shell
cd nats-streaming-server-v0.25.6
nats-streaming-server.exe
```

### 5. Запустите сервер

```shell
cd WB_Tech_L_0\cmd\app
go run main.go
```

### 6. Откройте приложение в своем браузере

Посетите сайт [http://localhost:8090](http://localhost:8090) в своем браузере.

### 7. Для публикации заказа

```shell
cd WB_Tech_L_0\tests
go run wb-order-service-test.go
```

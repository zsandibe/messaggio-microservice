# Messaggio

**A microservice on Go that accepts messages via REST API, stores them in PostgreSQL, and then sends them to Kafka for further processing. Processed messages are tagged. The service also provides an API to get statistics on processed messages.**


## Clone the project

```
$ git clone https://github.com/zsandibe/messaggio-microservice
$ cd messaggio-microservice
```

## Launch a project

```
$ make run
```

## Execute migrations

```
$ make migrate-up
$ make migrate-down
```

OR

```
$ make start
$ make stop
```

## SwaggerUI

```

localhost:8888/swagger/index.html

```

## Makefile tips

```

include .env
export

create-migrations:
	migrate create -ext sql -dir ./migrations -seq init_schema

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres


createdb: 
	docker exec -it postgres createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	sudo docker exec -it postgres dropdb  $(DB_NAME)


migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable" -verbose up 

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable" -verbose down



start:
	docker-compose up

stop:
	docker-compose down

all-containers-delete:
	docker ps -a -q | xargs -r docker rm -f


all-images-delete:
	docker images -q | xargs -r docker rmi -f



run:
	go run cmd/main.go


.PHONY: deps
deps:
	go mod tidy

.PHONY: swagger-init
swagger-init:
	@echo "Generate swagger gui"
	swag init -g  cmd/main.go



.PHONY: create-migrations postgres createdb dropdb migrate-up migrate-down start stop run 

```


## API server provides the following endpoints:
* `GET /api/v1/messages` - returns a messages list by filter
* `POST /api/v1/messages` - creates a messages  by body(content)
* `GET /api/v1/messages/{id}` - returns a message by id from query path
* `DELETE /api/v1/message/{id}` - deletes a message by id from path
* `GET /api/v1/stats` - returns a stats of messages
* `GET /api/v1/stats/{id}` - return stat of message by id from path


# .env file
## Kafka configuration

```
KAFKA_BROKER=
KAFKA_TOPIC=
KAFKA_GROUP_ID=

```

## Server configuration

```
SERVER_HOST=localhost
SERVER_PORT=8888
```

## Postgres configuration

```
DRIVER=
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
```


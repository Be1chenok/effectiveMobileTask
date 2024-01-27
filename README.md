# Effective Mobile Task

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)


## Clone the project

```
$ git clone https://github.com/Be1chenok/effectiveMobileTask
$ cd effectiveMobileTask
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

## Logs

Logs are saved in the logs folder:

* app.log - info and debug logs
* error.log - error logs

## API server provides the following endpoints:
* `GET /persons` - returns a list of persons (parameters can be used: gender, nationality, page, size)
* `GET /person/{id}` - returns person by id
* `POST /person` - adds a person

```
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich" <-optionally
}
```

* `PUT /person` - updates person by id

```
{
    "id": 1,
    "name": "Genadiy",
    "surname": "Shapkin",
    "patronymic": "Petrovich", <-optionally
    "age": 21,
    "gender":"male",
    "nationality":"RU"
}
```

* `DELETE /person/{id}` - deletes person by id

# .env file
## Enrichment configuration

```
AGE_URL=https://api.agify.io
GENDER_URL=https://api.genderize.io
NATIONALITY_URL=https://api.nationalize.io
```

## Server configuration

```
SERVER_HOST=server
SERVER_PORT=8080
REQUEST_TIME=5 <- in seconds
```

## Postgres configuration

```
PG_HOST=postgres
PG_PORT=5432
PG_USER=postgres
PG_PASS=postgres
PG_BASE=postgres
PG_SSL_MODE=disable
```
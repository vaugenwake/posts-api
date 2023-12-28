# Project posts-api

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

## Migrations

Migrations are supported by [Goose](https://github.com/pressly/goose) and can be run independently of the application.

### Start up
When the app is started in docker the migrations will be migrated up on first start
```BASH
make docker-run
```

### Migrate up (foreward)
Once you have made a new `*.sql` migration in the `./migrations` directory you can migrate the database forward by running
```BASH
make migrate-up
```
This will start the migrations docker container and execute `goose up`

### Migrate status
You can check the current state of the DB based on your migrations by using the status command
```BASH
make migrate-status
```
This will start the migrations docker container and execute `goose status`

### Migrate down (rollback)
If you have applied a migration and you would like to roll back the change you can use the down command. (This is not recommended as it might cause data deletion, rolling forward is always preferred)
```BASH
make migrate-down
```
This will start the migrations docker-container and execute `goose down` and step the DB version back by 1 increment
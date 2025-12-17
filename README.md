# Recipes API

Recipes Example API, This is a sample server celler server.

- [License](https://bramwork.com/terms/)
- [Support](https://bramworks.com/support)

## Server

- Host: localhost
- Port: 8080
- BasePath: /api/v1

## Documentation

- [Version 1.0.0](https://bramworks.com/resources/open-api/)

## Swagger

Create/update swagger model documentation
```sh
$ swag init
```

Running swagger server
```sh
$ swagger serve -F swagger ./docs/swagger.json
```

## Running
```sh
$ MONGO_URI="mongodb://<USERNAME>:<PASSWORD>@localhost:27017/test?authSource=admin" go run .
```
# GO GIN/MONGO API STARTER KIT

Starter kit application to get started with the development of REST API.

Stack:

- Go GIN
- MongoDB
- Swagger
- Docker

## Getting started

Build everything inside two docker containers: `mongo` and `app`.

This command will build a dev environment, mongo is used without user and password authentication.

``` bash
$ ./run.sh
```

A dev server runs at http://localhost:8080

Live relaod on code change is active.

### Generate OpenAPI documentation (Swagger)

``` bash
$ cd src
$ swag init
```

Then visit http://localhost:8080/swagger/index.html

## Utils

Use the `shell.sh` script to connect ssh into one of the two containers:

``` bash
$ ./shell.sh app
$ ./shell.sh mongo
```

## Production

Create a `.env.prod` file:

```
MONGO_INITDB_ROOT_USERNAME=root
MONGO_INITDB_ROOT_PASSWORD=mypass
MONGO_INITDB_DATABASE=mydb
MONGO_SUPERADMIN_EMAIL=admin
MONGO_SUPERADMIN_PASSWORD=password
APP_SETTINGS=settings.prod.json
```

Create a prod settings file `src/settings.prod.json`:

``` json
{
    "database": {
        "host": "mongodb://root:mypass@mongo",
        "port": 27017,
        "dbName": "mydb"
    },
    "jwt": {
        "secret": "Shfdjlkl$gfj!"
    }
}
```

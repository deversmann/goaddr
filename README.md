# goaddr

![CI](https://github.com/deversmann/goaddr/actions/workflows/ci.yml/badge.svg)

A Restful API (and db) for a super simple Address Book written in Golang.  Requires Go 1.15 or higher to build.

## Getting the app

### Go Binaries
The easiest way to access the app is via `go`. The following commands will download the executable into your go files and run the API server on port 8080 with an in-memory database:
```bash
go get github.com/deversmann/goaddr
goaddr
```

### Container
The latest build of the container can also be found at https://quay.io/repository/deversmann/goaddr and can be pulled using:
```bash
podman pull quay.io/deversmann/goaddr
```

### Building locally
Building the application requires go 1.15 or higher. Clone the repository and enter the directory and execute the following:
```bash
go mod tidy # downloads the dependencies

go run main.go # runs without building

#or

go build # builds the executable in the working directory
./goaddr # runs from the local directory

# or

go install # builds the executable and installs it in the go binaries directory
goaddr # assumes that go binaries is on your path
```

The project includes a `Containerfile` that can be used to build a container containing the application.  To build the container locally:
```bash
podman build --tag localhost/goaddr .
```
The resulting container could be run using:
```bash
podman run -d -rm \
  -e GOADDR_DBDIALECT="postgresql" \
  -e GOADDR_DBDSN="host=192.0.2.42 user=user password=pass dbname=db" \
  -e GOADDR_PORT=8081
  -p 8081:8081 --name goaddr localhost/goaddr
```

## Configuration

The following attributes can be configured by using environment variables:

| Configuration | ENV Var Name | Default | Notes |
|---|---|---|---|
| Database Dialect | `GOADDR_DBDIALECT` | `sqlite` | The type of database being connected to.  Currently only `sqlite` and `postgresql` are valid options. For `sqlite`, the service will create the specified db file if it doesn't exist. |
| Database DSN | `GOADDR_DBDSN` | `file::memory:?cache=shared`<br>(in-memory db) | The connect string for the database selected. See https://gorm.io/docs/connecting_to_the_database.html for examples |
| Web Service Port | `GOADDR_PORT` | `8080` | The port the web service will listen on |
| Logging Level | `GOADDR_LOGLEVEL` | No Default| There are only 2 log levels, DEBUG and INFO.  By default, only INFO is on.  Set to `DEBUG` to add DEBUG or to `NONE` to turn off all. |


## API Definition
| Call | Success | Failure |
|---|---|---|
| `POST /api/v1/contacts` | 201 created | 400 if JSON is not valid |
| `GET /api/v1/contacts` [*](#query-options) | 200 OK | 400 if query is not valid<br>404 if no contacts are found |
| `GET /api/v1/contacts/:id` | 200 OK | 404 if contact with id is not found |
| `PUT /api/v1/contacts/:id` | 202 accepted | 400 if JSON is not valid<br>404 if contact with id is not found |
| `DELETE /api/v1/contacts/:id` | 204 no content | 404 if contact with id is not found |

### Query options
The group GET request has multiple combinable ways of being limited:
| Type | Syntax | Notes |
|---|---|---|
| Filter | `<field>=<value>` | Does a case-insensitive search for the **value** in the named **field** |
| Sort | `sort_by=<field1>,-<field2>` | Prepend the field name with '-' for descending sort<br>Should be used for reliable pagination |
| Limit | `limit=<number>` | Only returns the first **number** results<br>Used in conjunction with **offset** to paginate results |
| Offset | `offset=<number>` | Discards the first **number** results<br>Used in conjunction with **limit** to paginate results |


## Return Values
If an entry is retrieved, created or modified, the api will return a JSON representation of that entry:

``` json
{
    "id": 1,
    "firstname": "Paul",
    "lastname": "Cormir",
    "address": "100 E. Davie St.",
    "city": "Raleigh",
    "state": "NC",
    "zipcode": "27601",
    "phone": "888-RED-HAT-1",
    "email": "pcormir@redhat.com"
}
```

If a collection is requested, the entries are returned in a JSON list:

```json
{
  "contacts": [
    {
      "id": 1,
      "firstname": "Paul",
      "lastname": "Cormir",
      "address": "100 E. Davie St.",
      "city": "Raleigh",
      "state": "NC",
      "zipcode": "27601",
      "phone": "888-RED-HAT-1",
      "email": "pcormir@redhat.com"
    },
    ...
  ]
}
```

If an error occurs, the HTTP error code will be accompanied by a message describing the failure:

```json
{
  "message": "Invalid JSON for contact",
  "status": 400
}
```

## Links

### Frameworks
- [GIN](https://gin-gonic.com) - Web framework for handling API requests and responses
- [GORM](https://gorm.io) - Used for database connections and mapping

### Other
The following pages were used significantly for inspiration and/or as references:
- https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/
- https://golang.org/doc/tutorial/web-service-gin
- https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
- https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin
- https://medium.com/wesionary-team/create-your-first-rest-api-with-golang-using-gin-gorm-and-mysql-d439bcc6f987



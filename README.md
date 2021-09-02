# goaddr

A Restful API (and db) for a super simple Address Book written in Golang.  Requires Go 1.16 or higher to be installed.

The following commands will download the executable into your go files and run the API server on port 8080:
``` bash
go get github.com/deversmann/goaddr
goaddr
```

The server will generate the db file the first time in the directory where you execute the command.  Subsequent executions in the same directory will reuse the same db file

## API Definition
| Call | Success | Failure |
|---|---|---|
| POST /api/v1/contacts/ | 201 created | 400 if JSON is not valid |
| GET /api/v1/contacts/ | 200 OK | 404 if no contacts are found |
| GET /api/v1/contacts/:id | 200 OK | 404 if contact with id is not found |
| PUT /api/v1/contacts/:id | 202 accepted | 400 if JSON is not valid<br>404 if contact with id is not found |
| DELETE /api/v1/contacts/:id | 204 no content | 404 if contact with id is not found |

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
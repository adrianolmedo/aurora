# Aurora server

Golang RESTfull API for test.

## Test with Postgres service:

```bash
git clone https://github.com/adrianolmedo/aurora.git
cp .env.example .env
docker-compose up -d --build postgres
```

**Join to `psql` and ingress the password `1234567a`:**

```bash
$ docker exec -it postgres /bin/sh
$ psql -U adrian -d aurora
```

**Install tables:**

```bash
\i tables.sql
\q
exit
```

**Start:**

```
docker-compose up -d --build app
```

## Endpoints:

**POST:** `/v1/users`

Body (JSON):

```json
{
    "name": "Becky"
}
```

Reponse (201 Created):

```json
{
    "message_ok": {
        "content": "user created"
    },
    "data": {
        "id": 11,
        "uuid": "def4eae8-7bf0-4102-9ba8-45fbccbd365e",
        "name": "Becky",
        "created": "2023-07-06T01:09:39.98926674-04:00"
    }
}
```

---

**GET:** `/v1/users/:id`

Reponse (200 OK):

```json
{
    "data": {
        "id": 11,
        "uuid": "def4eae8-7bf0-4102-9ba8-45fbccbd365e",
        "name": "Becky",
        "created": "2023-07-06T01:09:39.989267Z"
    }
}
```

---

**GET:** `/v1/users?limit=2&page=3`

Get all users but filtered by pagination.

Response (200 OK):

```json
{
    "data": [
        {
            "id": 6,
            "uuid": "6e79e75e-ed4a-4394-b7ab-5134386d8561",
            "name": "Daniela",
            "created": "2023-07-06T01:09:04.124354Z"
        },
        {
            "id": 7,
            "uuid": "ec6ef353-96ae-4d25-9721-7c48a7b64ee0",
            "name": "Rebeca",
            "created": "2023-07-06T01:09:09.739448Z"
        }
    ],
    "links": {
        "first": "/v1/users?limit=2&page=1&sort=created_at",
        "prev": "/v1/users?limit=2&page=2&sort=created_at",
        "next": "/v1/users?limit=2&page=4&sort=created_at",
        "last": "/v1/users?limit=2&page=5&sort=created_at"
    },
    "meta": {
        "limit": 2,
        "page": 3,
        "sort": "created_at",
        "total": 10,
        "total_pages": 5,
        "from_row": 7,
        "to_row": 8
    }
}
```

---

**DELETE:** `/v1/users/:id`

Response (200 OK):

```json
{
    "message_ok": {
        "content": "user deleted"
    }
}
```

## Run integration tests

Some storage layer related files have a build tag:

```go
// go:build integration
// +build integration
```

It also parses `package main` which calls `flag.Parse`, so all declared and visible flags will be parsed and available for the tests.

Example for run:

```bash
go test -v -tags integration -args -dbhost 127.0.0.1 -dbport 5432 -dbuser adrian -dbname testdb -dbpass 1234567b
```

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

**POST:** `/v1/signup`

Sing up users or create account. *First Name, Email and Password are fields required.*

Body (JSON):

```json
{
    "first_name": "John",
    "last_name": "Doe",
    "email": "jdoe@go.com",
    "password": "1234567b"
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

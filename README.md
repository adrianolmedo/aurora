# Aurora server

Golang RESTfull API for test.

## Test with Postgres service:

```bash
$ git clone https://github.com/adrianolmedo/aurora.git
$ cp .env.example .env
$ docker-compose up -d --build postgres
```

**Join to `psql` and ingress the password `1234567a`:**

```bash
$ docker exec -it postgres /bin/sh
$ psql -U adrian -d aurora
```

**Install tables:**

```bash
$ \i tables.sql
$ \q
$ exit
```

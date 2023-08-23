# codelibrary

## Development

Ensure Docker is installed, and start the project with `docker compose up`.

You can connect to the Postgres database through the `db` service:

```
docker compose exec db psql codelibrary postgres
```

The database schema will be created automatically on first run. You can update
the schema if needed by running the SQL file again.

```
docker compose exec db psql codelibrary postgres \
  -q -f /docker-entrypoint-initdb.d/codelibrary.sql
```

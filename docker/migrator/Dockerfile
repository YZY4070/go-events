# syntax=docker/dockerfile:1.7

FROM migrate/migrate:4

WORKDIR /migrator

COPY ../../migrations /migrations

ENV MIGRATIONS_PATH=/migrations

ENTRYPOINT ["sh", "-c", "migrate -path=${MIGRATIONS_PATH} -database ${DATABASE_URL} up"]

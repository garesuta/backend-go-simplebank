# Build stage
FROM golang:1.23.4 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# RUN apt-get update && apt-get install -y
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM golang:1.23.4
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .

# copy all migration to docker
COPY db/migration .db/migration

# inform docker port
EXPOSE 8080

# define default command
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh"]

#
FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o vocabsrv .

# ---

FROM alpine:latest
WORKDIR /app
COPY --from=build ./app/vocabsrv ./
COPY --from=build ./app/configs/config.docker.toml ./config.toml
EXPOSE 8080

CMD ["./vocabsrv"]

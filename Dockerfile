FROM golang:1.18 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/src/app/bin/app ./cmd/

FROM alpine:latest

WORKDIR /usr/src/app

COPY --chown=65534:65534 --from=build /usr/src/app/bin/app .
COPY --chown=65534:65534 --from=build /usr/src/app/configs/config.env .

USER 65534

EXPOSE 8000

CMD ["/usr/src/app/app"]
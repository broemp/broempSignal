FROM golang:1.19 AS BUILD

WORKDIR /app_build

COPY ./ ./
RUN go mod download
RUN go build -o /app/BroempSignal

FROM ubuntu:latest
COPY --from=BUILD /app/BroempSignal /app/BroempSignal
ENTRYPOINT  ["/app/BroempSignal"]
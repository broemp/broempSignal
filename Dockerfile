FROM golang:latest AS BUILD

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /broempSignal cmd/broempSignal.go

FROM alpine
WORKDIR /app
COPY --from=BUILD /broempSignal .
CMD ["/app/broempSignal"]

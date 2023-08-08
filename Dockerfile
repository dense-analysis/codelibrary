FROM golang:1.19

# TODO: Find a way to not do this for production builds.
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/cosmtrek/air@latest

WORKDIR /code

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build -o main ./cmd/codelibrary/main.go

EXPOSE 8080

CMD ["./main"]

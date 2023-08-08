# codelibrary

Install `swag` to generate swagger documentation.

```
go install github.com/swaggo/swag/cmd/swag@latest
```

```
swag init -d cmd/codelibrary/ -o internal/docs
```

Install `air` for hot code reloading.

```
go install github.com/cosmtrek/air@latest
```

You can run the project in development mode with `air`.

```
air
```

# morph-tracking-api

API for tracking morphing event on the ronin chain.

## Rebuild locally

> This assumes that you have a [correctly configured](https://golang.org/doc/install#testing) Go toolchain

Install the dependencies

```go 
go mod tidy
``` 

Create a `.env` file with corresponding values based on `.dist.env`

```
DB_USERNAME=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
```

Build the binary and run the executable

```bash
go build .
./morph-tracking-api
```
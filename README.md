# morph-tracking-api

API for tracking morphing event on the ronin chain.

## Usage

Use this [endpoint](https://morph.betteraxie.tech/) to use this API. This accepts the following parameters:

| Params | Value   | Description                                                                                                |
| ------ | ------- | ---------------------------------------------------------------------------------------------------------- |
| page   | Integer | Page of the result. Each page is an offset of page \* 100.                                                 |
| asc    | Boolean | Determines if results should be in ascending or descending order.                                          |
| byId   | Boolean | Determines if the order will be based on the Axie's ID or by the default setting. Default is by timestamp. |


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
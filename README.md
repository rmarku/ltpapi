# Golang Exercise


## Usage

```bash
go run .
```
Some configuration can be passed by config file ore environment variables. Configuration available are in [config.go](internal/config/config.go)

**Example with environment variables:**

```bash
LTPAPI_SERVER_PORT=:8888 go run .
```

## Testing

The application has some unit tests, but also integration tests that test some common flows.
For the unit test testify is used, for the integration test, testing containers is used.

To run the unit tests:
```bash
go test ./... -v 
```

To run the integration tests:
```bash
go test ./... -tags integration  -v
```

## Dockerized application

A Dockerfile is provided, it creates the image in two steps. The build step download the dependencies and build the executable.
The second step is used to run the application, it uses the builded executable from the previous step.

To build the image:
```bash
docker buildx build --load -t ltpapi .
```

To run the image:
```bash
docker run -p 8081:8081 ltpapi
```

Service can be tested with the following commands:
```bash
curl "http://localhost:8081/api/v1/ltp"
```
```bash
curl "http://localhost:8081/api/v1/ltp?pairs=BTC/USD,BTC/EUR"
```

## Design considerations

This project is designed with some hexagonal architecture approach. Here we have a
domain with the business rules and another layer for the infrastructure. This separation
is not in the folder structure, but in how the component are used.

We have `handlers` that are the primary adapters, the interface is the port and the only 
implementation the actual api endpoint. This handler call the business logic in the domain service. The service use two secondary adapters, a datasource (the call to the kraken API)
and a key-value store (to cache the kraken responses).

This design allows to make the code testeable (each adapter and the domain have unit test)
and also allows to change the implementation of each adapter without changing the code in
the domain. As an example of a posible improvement, a new keyvalue adapter can be created
to use redis instead of the simple map. This would allow to horizontally scale the
application and each replica to use the same cache.

The application follow 12 factor design principles, in particular the one to allow changing
configuration from the enviroment. This is done by using [viper](https://github.com/spf13/viper)
to read the configuration file and environment variables.

The actual functionality allow two run modes, one is with a ticker that refresh
the amount from *kraken api* every `ticker.timeout` to key the in memory caches,
and the other to refresh each time the cache expire. Each mode has its pros and cons:

- **Ticker**: The application is always up to date with the latest amount from kraken and
the response time always fast, as results came always from cache. The downside is that
the application is always doing call to Kraken API.
- **On demand**: The application is not always up to date and some queries to the api
can take a while. But this mode allow to reduce the amount of calls to kraken API.


### TODOs
- [ ] Fix all linters errors
- [ ] Add custom errors to domain
- [ ] Move applicaton constraction (main.go) to cmd directory
- [ ] Improve github actions workflow
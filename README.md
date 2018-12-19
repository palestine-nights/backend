[license]: ./LICENSE
[docs]: ./docs
[docker]: ./Dockerfile
[ci]: https://circleci.com/gh/palestine-nights/backend
[ci-badge]: https://circleci.com/gh/palestine-nights/backend.svg?style=svg

# Backend

[![Circle CI][ci-badge]][ci]

> REST API for table reservation

Created to avoid issue with CORS, which appears with axios in VueJS Apps.

## Development

Compile source code

```sh
$> go build -o main src/*.go
```

Run server

```
$> ./main
```

## Usage

Build and deploy using [docker][docker].

See [API documentation][docs] for more information.

Example with [httppie](https://httpie.org).

```sh
$> http GET http://localhost:8080/reservations
```

```json
[
    {
        "duration": 120,
        "email": "johndoe@example.com",
        "full_name": "John Doe",
        "guests": 5,
        "id": 1,
        "phone": "+380123456789",
        "state": "created",
        "table_id": 1,
        "time": "2019-11-25T23:50:00Z"
    }
]
```

## License
Project released under the terms of the MIT [license][license].

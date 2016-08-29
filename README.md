# env

Unix `env` utility for Windows users.
Supports minimal [POSIX functionality](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/env.html).

## Download

See [releases page](https://github.com/davidmz/env/releases).

## Get & build

```
go get https://github.com/davidmz/env
```

## Usage example

```
env GOOS=linux GOARCH=amd64 go build some/package
```


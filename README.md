# w84

w84 /weɪt fɔːr/ is a yet another clone of [wait-for-it.sh](https://github.com/vishnubob/wait-for-it).

## Usage

```
w84 -t 15s example.com:22 api.example.com:80 unix:/var/run/docker.sock
```

### Options

* `-v`, `-verbose` turns on the verbose mode which means
  letting the command print more about the errors while connectivity tests.
* `-t`, `-timeout` is duration of timeout for connectivity tests.

### Exit Code

* `0` means it seems OK for all the given sites connectivity.
* `1` means one or more sites didn't respond within the timeout.

Command returns non-zero and other than above exit code
if the command argument is empty (no sites are specified).

## Build

```
go test ./...
go build ./cmd/...
```

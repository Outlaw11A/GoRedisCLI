# GoRedisCLI

A toy example of a Redis CLI client written in Golang.

## Usage

`./gorediscli -h 127.0.0.1 -p 6379 -c "keys *"`

Host, port and command are completely optional. It will default to host 127.0.0.1, port 6379 and interactive mode)
# remote-log

A golang module for enabling remote logging on applications and servers. Start a log-broadcast server and connect from various clients using various protocols.

## Currently Supported Protocols

- tcp
- udp

## Logger Usage

see: `./examples/counting/main.go`

## Client Usage

run `./cmd/remote-log/main.go`

arguments: `<protocol>` `<host>` `<port>`
-  protocol: tcp or udp
-  host: the hosts ip address
-  port: the port where the server is running on

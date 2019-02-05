# TCP Server

honestbee take home test

## Build

It will create TCP Server, External API Server and Client at bin folder.

``` zsh
make
```

## Usage

### Server

Start TCP Server and External API Server.

``` zsh
./bin/server
./bin/external
```

### Client

Run the client to connect to tcp server.

``` zsh
# n is the number of client you want to connect to tcp server, default is 1.
./bin/client -n 10
```

### Observer

You can check TCP Server by browsing "http://localhost/state" or using curl.
The Script will check the url per second.

``` zsh
./curl.sh
```

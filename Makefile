all:
	# go build -o ./bin/server ./server/main.go
	# go build -o ./bin/client ./client/main.go
	go build -ldflags="-s -w" -o ./bin/server ./main.go ./httpsrv.go ./tcpsrv.go ./util.go
	go build -ldflags="-s -w" -o ./bin/client ./test/client/main.go
	go build -ldflags="-s -w" -o ./bin/external ./test/external/main.go
clean:
	rm ./bin/*
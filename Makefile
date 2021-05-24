export CLIENT_CMD='./dist/client --datadir /tmp/dropbox/client'
export SERVER_CMD='./dist/server --datadir /tmp/dropbox/server'

.PHONY: proto client server

build: client server

client:
	go build -o dist/client cmd/client/*.go

server:
	go build -o dist/server cmd/server/*.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    proto/gobox.proto

test: build
	rm -rf /tmp/dropbox
	mkdir -p /tmp/dropbox/client /tmp/dropbox/server
	pytest -vv -s .

test-ci: build
	mkdir -p /tmp/dropbox/client /tmp/dropbox/server
	pytest -q -rapP

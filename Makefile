# export CLIENT_CMD='python3 -c "import test_homework as th; th.client()" -- dist/client --datadir /tmp/dropbox/client'
# export SERVER_CMD='python3 -c "import test_homework as th; th.server()" -- dist/server --datadir /tmp/dropbox/server'
# export CLIENT_CMD='python3 -c "import test_homework as th; th.client()" -- ls'
# export SERVER_CMD='python3 -c "import test_homework as th; th.server()" -- ls ~'
export CLIENT_CMD='./dist/client --datadir /tmp/dropbox/client'
export SERVER_CMD='./dist/server --datadir /tmp/dropbox/server'

.PHONY: proto client server

build: proto client server

client:
	go build -o dist/client client/*.go

server:
	go build -o dist/server server/*.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    proto/gobox.proto

test: build
	rm -rf /tmp/dropbox
	mkdir -p /tmp/dropbox/client /tmp/dropbox/server
	#pytest -vv -s .
	pytest -vv -s . -k 'test_add_single_file'




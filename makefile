all: get_all build_server build_client

build_server:
	cd ./server && go build -o ../serverapp

build_client:
	cd ./client && go build -o ../clientapp

get_all:
	go get "github.com/karetskiiVO/NATpunching/natpunch"
	go get "github.com/jessevdk/go-flags"